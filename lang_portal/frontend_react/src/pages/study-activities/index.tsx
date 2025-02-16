import { FC, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { studyActivitiesService } from '../../services/study-activities'
import type { StudyActivity } from '../../types/study-activities'
import { getActivityThumbnailUrl } from '../../utils/image-utils'

const StudyActivitiesPage: FC = () => {
  const [activities, setActivities] = useState<StudyActivity[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchActivities = async () => {
      try {
        const response = await studyActivitiesService.getAll()
        setActivities(response.data.items)
      } catch (err) {
        console.error('Error loading study activities:', err)
        setError('Failed to load study activities')
      } finally {
        setIsLoading(false)
      }
    }

    fetchActivities()
  }, [])

  if (isLoading) {
    return (
      <div className="p-8">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-white rounded-lg shadow-md p-6">
                <div className="w-full h-48 bg-gray-200 rounded-md mb-4"></div>
                <div className="h-6 bg-gray-200 rounded w-3/4 mb-4"></div>
                <div className="flex justify-between">
                  <div className="h-10 bg-gray-200 rounded w-24"></div>
                  <div className="h-10 bg-gray-200 rounded w-24"></div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="p-8">
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      </div>
    )
  }

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-6">Study Activities</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {activities.map(activity => (
          <div key={activity.id} className="bg-white rounded-lg shadow-md p-6">
            <img 
              src={getActivityThumbnailUrl(activity.name, activity.thumbnail_url)}
              alt={activity.name}
              className="w-full h-48 object-cover rounded-md mb-4"
              onError={(e) => {
                const img = e.target as HTMLImageElement;
                img.src = '/placeholder.png';
              }}
            />
            <h2 className="text-xl font-semibold mb-4">{activity.name}</h2>
            {activity.description && (
              <p className="text-gray-600 mb-4">{activity.description}</p>
            )}
            <div className="flex justify-between">
              <Link 
                to={`/study_activities/${activity.id}`}
                className="px-4 py-2 bg-gray-100 rounded-md hover:bg-gray-200"
              >
                View Sessions
              </Link>
              <Link 
                to={activity.launch_url}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                Launch
              </Link>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

export default StudyActivitiesPage
