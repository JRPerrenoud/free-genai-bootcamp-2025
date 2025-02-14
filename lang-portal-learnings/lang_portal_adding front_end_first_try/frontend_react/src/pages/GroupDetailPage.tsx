import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { PaginatedTable, type Column } from '@/components/common/PaginatedTable';
import type { Word, WordGroup } from '@/types/words';
import { groupService } from '@/services/groups';

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
    header: 'Correct',
    accessorKey: 'correct_count',
  },
  {
    header: 'Wrong',
    accessorKey: 'wrong_count',
  },
];

export default function GroupDetailPage() {
  const { partOfSpeech } = useParams();
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [group, setGroup] = useState<WordGroup | null>(null);
  const [words, setWords] = useState<Word[]>([]);
  const [totalPages, setTotalPages] = useState(1);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchGroupDetails = async () => {
      if (!partOfSpeech) return;
      
      try {
        const [groupResponse, wordsResponse] = await Promise.all([
          groupService.getGroupByPartOfSpeech(partOfSpeech),
          groupService.getGroupWords(partOfSpeech, page),
        ]);

        setGroup(groupResponse.data);
        setWords(wordsResponse.data.items);
        setTotalPages(wordsResponse.data.total_pages);
      } catch (error) {
        console.error('Failed to fetch group details:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchGroupDetails();
  }, [partOfSpeech, page]);

  if (isLoading) {
    return <div className="h-full w-full flex items-center justify-center">Loading...</div>;
  }

  if (!group) {
    return <div>Group not found</div>;
  }

  return (
    <div className="h-full flex flex-col space-y-6 p-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold">{group.name}</h1>
          <p className="text-gray-600">Total Words: {group.word_count}</p>
        </div>
        <button
          onClick={() => navigate('/words/new')}
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-md"
        >
          Add Word
        </button>
      </div>

      <div className="flex-1">
        <PaginatedTable
          data={words}
          columns={columns}
          currentPage={page}
          totalPages={totalPages}
          onPageChange={setPage}
          onRowClick={(row) => navigate(`/words/${row.id}`)}
        />
      </div>
    </div>
  );
}
