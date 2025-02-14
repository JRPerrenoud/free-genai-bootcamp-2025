import { HTMLAttributes } from 'react';
import { cn } from '../../lib/utils';
import { Button } from './Button';
import { type ReactNode } from 'react';

export interface Column<T = any> {
  header: string;
  accessorKey: keyof T | string;
  cell?: ({ row, table }: { row: { original: T }; table: { options: { meta?: any } } }) => ReactNode;
  className?: string;
}

interface PaginatedTableProps<T> extends HTMLAttributes<HTMLDivElement> {
  data: T[];
  columns: Column<T>[];
  page: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  meta?: Record<string, any>;
}

export function PaginatedTable<T>({
  data,
  columns,
  page,
  totalPages,
  onPageChange,
  meta,
  className,
  ...props
}: PaginatedTableProps<T>) {
  const getCellContent = (row: T, column: Column<T>) => {
    if (column.cell) {
      return column.cell({ row: { original: row }, table: { options: { meta } } });
    }
    if (typeof column.accessorKey === 'function') {
      return column.accessorKey(row);
    }
    return row[column.accessorKey] as string;
  };

  return (
    <div className={cn('h-full flex flex-col', className)} {...props}>
      <div className="flex-1 overflow-auto">
        <table className="w-full border-collapse">
          <thead className="sticky top-0 bg-gray-50">
            <tr className="border-b">
              {columns.map((column, index) => (
                <th
                  key={index}
                  className={cn(
                    'px-6 py-4 text-left text-sm font-semibold text-gray-900',
                    column.className
                  )}
                >
                  {column.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {data.map((row, rowIndex) => (
              <tr
                key={rowIndex}
                className="border-b hover:bg-gray-50"
              >
                {columns.map((column, colIndex) => (
                  <td
                    key={colIndex}
                    className={cn(
                      'px-6 py-4 text-sm text-gray-500',
                      column.className
                    )}
                  >
                    {getCellContent(row, column)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      
      <div className="flex items-center justify-between px-6 py-4 bg-white border-t">
        <div className="text-sm text-gray-500">
          Page {page} of {totalPages}
        </div>
        <div className="flex gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(page - 1)}
            disabled={page <= 1}
          >
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(page + 1)}
            disabled={page >= totalPages}
          >
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}
