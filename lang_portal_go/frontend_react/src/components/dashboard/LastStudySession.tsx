import { FC } from 'react'
import { Link } from 'react-router-dom'
import type { LastStudySession as LastStudySessionType } from '../../types/dashboard'

interface LastStudySessionProps {
  data: LastStudySessionType | null
  isLoading: boolean
}

const LastStudySession: FC<LastStudySessionProps> = ({ data, isLoading }) => {
  if (isLoading) {
    return (
      <div className="p-6 bg-white rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Last Study Session</h2>
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-3/4 mb-4"></div>
          <div className="h-4 bg-gray-200 rounded w-1/2"></div>
        </div>
      </div>
    )
  }

  if (!data) {
    return (
      <div className="p-6 bg-white rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Last Study Session</h2>
        <p className="text-gray-600">No study sessions yet</p>
      </div>
    )
  }

  const formattedDate = new Date(data.created_at).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })

  return (
    <div className="p-6 bg-white rounded-lg shadow">
      <h2 className="text-xl font-semibold mb-4">Last Study Session</h2>
      <div className="space-y-4">
        <div>
          <p className="text-gray-600">Group</p>
          <Link to={`/groups/${data.group_id}`} className="text-blue-600 hover:underline">
            {data.group.name}
          </Link>
        </div>
        <div>
          <p className="text-gray-600">Activity</p>
          <Link to={`/study_activities/${data.study_activity_id}`} className="text-blue-600 hover:underline">
            {data.activity.name}
          </Link>
        </div>
        <div>
          <p className="text-gray-600">Date</p>
          <p>{formattedDate}</p>
        </div>
        <div>
          <p className="text-gray-600">Items Reviewed</p>
          <p>{data.review_items_count}</p>
        </div>
      </div>
    </div>
  )
}

export default LastStudySession
