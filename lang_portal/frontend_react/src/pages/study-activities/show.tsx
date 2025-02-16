import { FC, useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { studyActivitiesService } from '../../services/study-activities'
import { studySessionsService } from '../../services/study-sessions'
import type { StudyActivity } from '../../types/study-activities'
import type { StudySession } from '../../types/study-sessions'

const StudyActivityShowPage: FC = () => {
  const { id } = useParams<{ id: string }>()
  const [activity, setActivity] = useState<StudyActivity | null>(null)
  const [sessions, setSessions] = useState<StudySession[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchData = async () => {
      if (!id) return
      
      try {
        const [activityRes, sessionsRes] = await Promise.all([
          studyActivitiesService.getById(parseInt(id)),
          studySessionsService.getByActivityId(parseInt(id))
        ])

        setActivity(activityRes.data)
        setSessions(sessionsRes.data.items)
      } catch (err) {
        console.error('Error loading study activity data:', err)
        setError('Failed to load study activity data')
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [id])

  if (isLoading) {
    return (
      <div className="p-8">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
          <div className="h-48 bg-gray-200 rounded-lg mb-6"></div>
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-16 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  if (error || !activity) {
    return (
      <div className="p-8">
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error || 'Study activity not found'}
        </div>
      </div>
    )
  }

  return (
    <div className="p-8">
      <div className="flex items-start justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold mb-2">{activity.name}</h1>
          {activity.description && (
            <p className="text-gray-600 mb-4">{activity.description}</p>
          )}
        </div>
        <Link
          to={activity.launch_url}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          Launch Activity
        </Link>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-semibold mb-4">Study Sessions</h2>
        
        <div className="overflow-x-auto">
          <table className="min-w-full">
            <thead>
              <tr className="bg-gray-50">
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Group</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Start Time</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">End Time</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Review Items</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {sessions.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-4 text-center text-gray-500">
                    No study sessions found
                  </td>
                </tr>
              ) : (
                sessions.map((session) => (
                  <tr key={session.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <Link 
                        to={`/sessions/${session.id}`}
                        className="text-blue-600 hover:text-blue-800"
                      >
                        {session.id}
                      </Link>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">{session.group_name}</td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {new Date(session.start_time).toLocaleString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {new Date(session.end_time).toLocaleString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">{session.review_items_count}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}

export default StudyActivityShowPage
