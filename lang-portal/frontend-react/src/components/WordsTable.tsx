import { ChevronUp, ChevronDown } from 'lucide-react'
import { Word } from '../services/api'

export type WordSortKey = 'spanish' | 'english' | 'correct_count' | 'wrong_count'

interface WordsTableProps {
  words: Word[]
  sortKey: WordSortKey
  sortDirection: 'asc' | 'desc'
  onSort: (key: WordSortKey) => void
}

export default function WordsTable({ words, sortKey, sortDirection, onSort }: WordsTableProps) {
  return (
    <div className="overflow-x-auto bg-white dark:bg-gray-800 rounded-lg shadow">
      <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead className="bg-gray-50 dark:bg-gray-900">
          <tr>
            {(['spanish', 'english', 'correct_count', 'wrong_count'] as const).map((key) => (
              <th
                key={key}
                scope="col"
                className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 dark:text-gray-100 sm:pl-6"
                onClick={() => onSort(key)}
                style={{ cursor: 'pointer' }}
              >
                <div className="flex items-center">
                  <span>
                    {key === 'spanish' && 'Spanish'}
                    {key === 'english' && 'English'}
                    {key === 'correct_count' && 'Correct'}
                    {key === 'wrong_count' && 'Wrong'}
                  </span>
                  {sortKey === key && (
                    sortDirection === 'asc' ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />
                  )}
                </div>
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200 dark:bg-gray-800 dark:divide-gray-700">
          {words.map((word) => (
            <tr key={word.id} className="hover:bg-gray-50 dark:hover:bg-gray-700">
              <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 dark:text-gray-100 sm:pl-6">
                {word.spanish}
              </td>
              <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500 dark:text-gray-300">
                {word.english}
              </td>
              <td className="px-3 py-4 whitespace-nowrap text-sm text-green-500 dark:text-green-400">
                {word.correct_count}
              </td>
              <td className="px-3 py-4 whitespace-nowrap text-sm text-red-500 dark:text-red-400">
                {word.wrong_count}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
