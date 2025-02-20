import streamlit as st
from typing import Dict, Callable
from question_storage import QuestionStorage

def render_question_list(
    storage: QuestionStorage,
    on_question_load: Callable[[Dict], None],
    container_width: int = None
) -> None:
    """
    Render the list of saved questions in a sidebar or container
    
    Args:
        storage: QuestionStorage instance
        on_question_load: Callback function when a question is loaded
        container_width: Optional width for the container
    """
    questions = storage.get_all_questions()
    
    if not questions:
        st.info("No saved questions yet. Generate and save some questions to see them here!")
        return
        
    # Group questions by practice type
    questions_by_type = {}
    for question in questions:
        practice_type = question.get("practice_type", "Other")
        if practice_type not in questions_by_type:
            questions_by_type[practice_type] = []
        questions_by_type[practice_type].append(question)
    
    # Display questions grouped by type
    for practice_type, type_questions in questions_by_type.items():
        with st.expander(f"📚 {practice_type} ({len(type_questions)})", expanded=False):
            for question in sorted(type_questions, key=lambda x: x.get("timestamp", ""), reverse=True):
                col1, col2 = st.columns([3, 1])
                
                # Show question preview
                with col1:
                    timestamp = question.get("timestamp", "Unknown time")
                    # Convert ISO timestamp to more readable format
                    try:
                        from datetime import datetime
                        dt = datetime.fromisoformat(timestamp)
                        timestamp = dt.strftime("%Y-%m-%d %H:%M")
                    except:
                        pass
                        
                    st.markdown(f"**{timestamp}**")
                    preview = question.get("introduction", "")[:100]
                    if len(question.get("introduction", "")) > 100:
                        preview += "..."
                    st.write(preview)
                
                # Action buttons
                with col2:
                    # Load button
                    if st.button("Load", key=f"load_q_{question['id']}"):
                        on_question_load(question)
                    
                    # Delete button
                    if st.button("Delete", key=f"del_q_{question['id']}"):
                        if storage.delete_question(question["id"]):
                            st.rerun()
                
                st.markdown("---")
