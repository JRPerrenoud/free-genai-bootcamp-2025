import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { studyActivitiesService } from '../../services/study-activities'
import type { StudyActivity, StudySession } from '../../types/study-activities'
import { StudySessionsList } from '../../components/study/StudySessionsList'
import { getActivityThumbnailUrl } from '../../utils/image-utils'

export default function StudyActivityShowPage() {
  const { id } = useParams<{ id: string }>()
  const [activity, setActivity] = useState<StudyActivity | null>(null)
  const [sessions, setSessions] = useState<StudySession[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchData = async () => {
      if (!id) return
      
      try {
        // First get the activity details
        const activityRes = await studyActivitiesService.getById(parseInt(id))
        setActivity(activityRes.data)

        // Then try to get the sessions
        try {
          const sessionsRes = await studyActivitiesService.getSessions(parseInt(id))
          setSessions(sessionsRes.data.items || [])
        } catch (err) {
          console.error('Error loading sessions:', err)
          setSessions([]) // Set empty sessions but don't show error
        }
      } catch (err) {
        console.error('Error loading study activity:', err)
        setError('Failed to load study activity')
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
          <img 
            src={getActivityThumbnailUrl(activity.name, activity.thumbnail_url)}
            alt={activity.name}
            className="w-48 h-48 object-cover rounded-lg shadow-md mb-4"
            onError={(e) => {
              const img = e.target as HTMLImageElement;
              img.src = '/placeholder.png';
            }}
          />
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
        <StudySessionsList sessions={sessions} activityName={activity.name} />
      </div>
    </div>
  )
}
