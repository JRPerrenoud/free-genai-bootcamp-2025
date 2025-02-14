import { useState, useEffect } from 'react';
import { PaginatedTable, type Column } from '@/components/common/PaginatedTable';
import { WordFormModal } from '@/components/words/WordFormModal';
import type { Word } from '@/types/words';
import { wordService } from '@/services/words';

const columns: Column<Word>[] = [
  {
    header: 'Spanish',
    accessorKey: 'spanish',
  },
  {
    header: 'English',
    accessorKey: 'english',
  },
  {
    header: 'Part of Speech',
    accessorKey: 'part_of_speech',
  },
  {
    header: 'Actions',
    accessorKey: 'id',
    cell: ({ row, table }) => {
      const { handleEdit, handleDelete } = table.options.meta as {
        handleEdit: (word: Word) => void;
        handleDelete: (id: number) => void;
      };
      return (
        <div className="flex space-x-2">
          <button
            onClick={() => handleEdit(row.original)}
            className="text-blue-600 hover:text-blue-800"
          >
            Edit
          </button>
          <button
            onClick={() => handleDelete(row.original.id)}
            className="text-red-600 hover:text-red-800"
          >
            Delete
          </button>
        </div>
      );
    },
  },
];

export default function WordsPage() {
  const [page, setPage] = useState(1);
  const [words, setWords] = useState<Word[]>([]);
  const [totalPages, setTotalPages] = useState(1);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedWord, setSelectedWord] = useState<Word | undefined>();

  const fetchWords = async () => {
    try {
      console.log('Fetching words, page:', page);
      const response = await wordService.getWords(page);
      console.log('Fetch response:', JSON.stringify(response, null, 2));
      
      if (!response.success || !response.data) {
        console.error('Failed to fetch words - invalid response:', response);
        return;
      }

      setWords(response.data.items);
      setTotalPages(response.data.total_pages);
    } catch (error) {
      console.error('Failed to fetch words:', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchWords();
  }, [page]);

  const handleAddWord = async (word: Omit<Word, 'id'>) => {
    try {
      console.log('Adding word:', word);
      const response = await wordService.createWord(word);
      console.log('Add word response:', JSON.stringify(response, null, 2));
      
      if (!response.success || !response.data) {
        console.error('Failed to add word - invalid response:', response);
        return false;
      }

      setPage(1); // Reset to first page
      await fetchWords();
      return true;
    } catch (error) {
      console.error('Failed to add word:', error);
      return false;
    }
  };

  const handleEditWord = async (id: number, word: Omit<Word, 'id'>) => {
    try {
      console.log('Editing word:', word);
      const response = await wordService.updateWord(id, word);
      console.log('Edit word response:', JSON.stringify(response, null, 2));
      
      if (!response.success || !response.data) {
        console.error('Failed to update word - invalid response:', response);
        return false;
      }

      await fetchWords();
      return true;
    } catch (error) {
      console.error('Failed to update word:', error);
      return false;
    }
  };

  const handleDeleteWord = async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this word?')) return;
    
    try {
      console.log('Deleting word:', id);
      const response = await wordService.deleteWord(id);
      console.log('Delete word response:', JSON.stringify(response, null, 2));
      
      if (!response.success) {
        console.error('Failed to delete word - invalid response:', response);
        return false;
      }

      await fetchWords();
      return true;
    } catch (error) {
      console.error('Failed to delete word:', error);
      return false;
    }
  };

  const handleEdit = (word: Word) => {
    setSelectedWord(word);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedWord(undefined);
  };

  if (isLoading) {
    return <div className="h-full w-full flex items-center justify-center">Loading...</div>;
  }

  return (
    <div className="h-full flex flex-col space-y-6 p-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-semibold text-gray-900">Words</h1>
        <button
          onClick={() => setIsModalOpen(true)}
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md"
        >
          Add Word
        </button>
      </div>

      <div className="flex-1 bg-white shadow rounded-lg overflow-hidden min-h-0">
        <PaginatedTable
          data={words}
          columns={columns}
          page={page}
          totalPages={totalPages}
          onPageChange={setPage}
          meta={{
            handleEdit,
            handleDelete: handleDeleteWord,
          }}
        />
      </div>

      <WordFormModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        onSubmit={selectedWord ? 
          (word) => handleEditWord(selectedWord.id, word) : 
          handleAddWord}
        initialData={selectedWord}
      />
    </div>
  );
}
