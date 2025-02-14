import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { PaginatedTable, type Column } from '@/components/common/PaginatedTable';
import type { WordGroup } from '@/types/words';
import { groupService } from '@/services/groups';

const columns: Column<WordGroup>[] = [
  {
    header: 'Part of Speech',
    accessorKey: 'name',
  },
  {
    header: 'Word Count',
    accessorKey: 'word_count',
  },
];

export default function GroupsPage() {
  const navigate = useNavigate();
  const [groups, setGroups] = useState<WordGroup[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchGroups = async () => {
      try {
        const response = await groupService.getGroups();
        setGroups(response.data.items);
      } catch (error) {
        console.error('Failed to fetch groups:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchGroups();
  }, []);

  if (isLoading) {
    return <div className="h-full w-full flex items-center justify-center">Loading...</div>;
  }

  return (
    <div className="h-full flex flex-col space-y-6 p-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Word Groups by Part of Speech</h1>
      </div>

      <div className="flex-1">
        <PaginatedTable
          data={groups}
          columns={columns}
          totalPages={1}
          currentPage={1}
          onPageChange={() => {}}
          onRowClick={(row) => navigate(`/groups/${row.name}/words`)}
        />
      </div>
    </div>
  );
}
