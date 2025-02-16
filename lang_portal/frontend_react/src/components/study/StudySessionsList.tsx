import { Link } from 'react-router-dom'
import { StudySession } from '../../types/study-activities'

interface StudySessionsListProps {
  sessions: StudySession[]
  activityName: string
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

export function StudySessionsList({ sessions, activityName }: StudySessionsListProps) {
  if (sessions.length === 0) {
    return (
      <div className="text-center py-4 text-gray-500">
        No study sessions found.
      </div>
    )
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              ID
            </th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Activity
            </th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Group
            </th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Start Time
            </th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Review Items
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {sessions.map((session) => (
            <tr key={session.id}>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600">
                <Link to={`/sessions/${session.id}`}>
                  {session.id}
                </Link>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {activityName}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {session.group.name}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {formatDate(session.created_at)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {session.review_items_count}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
