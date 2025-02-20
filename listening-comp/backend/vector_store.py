from typing import List, Dict, Optional, Any
import chromadb
from chromadb.config import Settings
import boto3
import json
from dataclasses import dataclass
import os

# Import DELEQuestion and TranscriptStructurer from structured_data
from structured_data import DELEQuestion, TranscriptStructurer

class BedrockEmbeddingFunction:
    def __init__(self, region_name: str = "us-east-2"):
        """Initialize Bedrock client"""
        self.bedrock = boto3.client(
            service_name='bedrock-runtime',
            region_name=region_name
        )
        self.model_id = "amazon.titan-embed-text-v2:0"

    def __call__(self, input: List[str]) -> List[List[float]]:
        """Generate embeddings for a list of texts using Amazon Bedrock"""
        embeddings = []
        
        for text in input:
            try:
                body = json.dumps({
                    "inputText": text
                })
                response = self.bedrock.invoke_model(
                    modelId=self.model_id,
                    body=body
                )
                response_body = json.loads(response['body'].read())
                embedding = response_body['embedding']
                embeddings.append(embedding)
            except Exception as e:
                print(f"Error generating embedding: {str(e)}")
                # Return a zero vector as fallback - now 1024 dimensions for v2 model
                embeddings.append([0.0] * 1024)
                
        return embeddings

def parse_questions_from_text(text: str) -> List[DELEQuestion]:
    """Parse XML-like formatted questions into DELEQuestion objects"""
    questions = []
    # Split by question tags
    question_blocks = text.split('<question')[1:]
    
    for block in question_blocks:
        try:
            # Remove closing tag and split into lines
            block = block.split('</question>')[0]
            parts = block.split('\n')
            intro = None
            conv = None
            quest = None
            current_section = None
            
            for part in parts:
                part = part.strip()
                if part.startswith('Introduction:'):
                    current_section = 'intro'
                elif part.startswith('Conversation:'):
                    current_section = 'conv'
                elif part.startswith('Question:'):
                    current_section = 'quest'
                elif part and current_section:
                    if current_section == 'intro':
                        intro = part
                    elif current_section == 'conv':
                        if conv:
                            conv += '\n' + part
                        else:
                            conv = part
                    elif current_section == 'quest':
                        quest = part
            
            if intro and conv and quest:
                questions.append(DELEQuestion(
                    introduction=intro,
                    conversation=conv,
                    question=quest
                ))
        except Exception as e:
            print(f"Error parsing question block: {str(e)}")
            continue
    
    return questions

class QuestionVectorStore:
    def __init__(self, persist_directory: str = "vectorstore"):
        """Initialize ChromaDB client with persistence"""
        self.persist_directory = persist_directory
        self.client = chromadb.PersistentClient(path=persist_directory)
        # Use Amazon Bedrock's Titan embedding model
        self.embedding_function = BedrockEmbeddingFunction()
        self.collection = self._get_or_create_collection()

    def _get_or_create_collection(self):
        """Get or create the collection, resetting if dimensions mismatch"""
        try:
            # Try to get existing collection
            try:
                collection = self.client.get_collection(
                    name="dele_questions",
                    embedding_function=self.embedding_function
                )
                print("Found existing collection")
                return collection
            except Exception as e:
                if "does not exist" in str(e):
                    print("Collection does not exist, creating new one...")
                    return self.client.create_collection(
                        name="dele_questions",
                        embedding_function=self.embedding_function
                    )
                raise
                
        except Exception as e:
            if "dimension" in str(e).lower():
                print("Dimension mismatch detected. Resetting collection...")
                self.reset_store()
                return self.client.create_collection(
                    name="dele_questions",
                    embedding_function=self.embedding_function
                )
            else:
                print(f"Unexpected error with collection: {str(e)}")
                raise

    def reset_store(self):
        """Reset the vector store by deleting and recreating it"""
        try:
            self.client.delete_collection("dele_questions")
            print("Deleted existing collection")
        except Exception as e:
            print(f"Error deleting collection (may not exist): {str(e)}")

    def add_question(self, question: DELEQuestion) -> str:
        """
        Add a question to the vector store
        Returns the ID of the added question
        """
        # Create a combined text for embedding
        combined_text = f"{question.introduction}\n{question.conversation}\n{question.question}"
        
        # Generate a unique ID
        question_id = str(hash(combined_text))
        
        # Add to collection with metadata containing the separate parts
        self.collection.add(
            documents=[combined_text],
            metadatas=[{
                "type": "original",
                "introduction": question.introduction,
                "conversation": question.conversation,
                "question": question.question
            }],
            ids=[question_id]
        )
        return question_id

    def find_similar_questions(self, query: str, n_results: int = 5) -> Dict:
        """Find similar questions based on semantic search"""
        try:
            # Check if collection is empty
            all_ids = self.collection.get()['ids']
            if not all_ids:
                print("Vector store is empty - no questions found")
                return None
                
            results = self.collection.query(
                query_texts=[query],
                n_results=min(n_results, len(all_ids)),  # Don't request more results than we have
                include=["metadatas", "distances"]
            )
            
            # Check if we got any results
            if not results['ids'][0]:
                print("No similar questions found for query")
                return None
            
            # Format results for better readability
            formatted_results = {
                'ids': results['ids'][0],
                'documents': [],
                'distances': results['distances'][0]
            }
            
            for metadata in results['metadatas'][0]:
                formatted_doc = (
                    f"Introduction:\n{metadata['introduction']}\n\n"
                    f"Conversation:\n{metadata['conversation']}\n\n"
                    f"Question:\n{metadata['question']}"
                )
                formatted_results['documents'].append(formatted_doc)
            
            return formatted_results
            
        except Exception as e:
            print(f"Error searching vector store: {str(e)}")
            return None

    def bulk_add_questions(self, questions: List[DELEQuestion]) -> List[str]:
        """Add multiple questions at once and return their IDs"""
        documents = []
        metadatas = []
        ids = []
        
        for question in questions:
            combined_text = f"{question.introduction}\n{question.conversation}\n{question.question}"
            question_id = str(hash(combined_text))
            
            documents.append(combined_text)
            metadatas.append({
                "type": "original",
                "introduction": question.introduction,
                "conversation": question.conversation,
                "question": question.question
            })
            ids.append(question_id)
        
        self.collection.add(
            documents=documents,
            metadatas=metadatas,
            ids=ids
        )
        return ids

    def get_question_by_id(self, question_id: str) -> Optional[DELEQuestion]:
        """Retrieve a specific question by ID"""
        result = self.collection.get(ids=[question_id])
        if result and result['documents']:
            return self.parse_document_to_question(result['documents'][0], result['metadatas'][0][0])
        return None

    def parse_document_to_question(self, document: str, metadata: Dict) -> Optional[DELEQuestion]:
        """
        Parse a document string back into a DELEQuestion object
        Format: introduction\nconversation\nquestion
        """
        try:
            return DELEQuestion(
                introduction=metadata['introduction'],
                conversation=metadata['conversation'],
                question=metadata['question']
            )
        except Exception as e:
            print(f"Error parsing document: {str(e)}")
        return None


if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser(description='Process YouTube transcript and store questions')
    parser.add_argument('transcript_path', help='Path to the transcript file')
    args = parser.parse_args()

    # Initialize vector store
    vector_store = QuestionVectorStore()
    
    # Get structured questions from transcript
    structurer = TranscriptStructurer()
    structured_text = structurer.structure_transcript(args.transcript_path)
    
    if not structured_text:
        print("Error: Failed to structure transcript")
        exit(1)
    
    # Parse questions
    questions = parse_questions_from_text(structured_text)
    print(f"Found {len(questions)} questions in transcript")
    
    # Store questions in vector store
    question_ids = vector_store.bulk_add_questions(questions)
    print(f"Successfully stored {len(question_ids)} questions")
    
    # Test similarity search with first question
    if questions:
        print("\nTesting similarity search with first question...")
        search_text = questions[0].introduction
        similar = vector_store.find_similar_questions(search_text)
        print(f"\nSearch query: {search_text}")
        print("\nSimilar questions found:")
        for doc, distance in zip(similar['documents'], similar['distances']):
            print(f"\nDocument (distance: {distance:.3f}):")
            print(doc)
