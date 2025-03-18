import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button } from "@/components/ui/button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useNavigation } from '@/context/NavigationContext'
import { createStudySession } from '@/services/api'

type Group = {
  id: number
  name: string
}

type StudyActivity = {
  id: number
  title: string
  launch_url: string
  preview_url: string
}

type LaunchData = {
  activity: StudyActivity
  groups: Group[]
}

export default function StudyActivityLaunch() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { setCurrentStudyActivity } = useNavigation()
  const [launchData, setLaunchData] = useState<LaunchData | null>(null)
  const [selectedGroup, setSelectedGroup] = useState<string>('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetch(`http://localhost:5000/api/study_activities/${id}/launch`)
      .then(response => {
        if (!response.ok) throw new Error('Failed to fetch launch data')
        return response.json()
      })
      .then(data => {
        setLaunchData(data)
        setCurrentStudyActivity(data.activity)
        setLoading(false)
      })
      .catch(err => {
        setError(err.message)
        setLoading(false)
      })
  }, [id, setCurrentStudyActivity])

  const handleLaunch = async () => {
    if (!selectedGroup) {
      setError('Please select a word group')
      return
    }

    // For typing tutor, we'll create a study session and open it in a new window with the group ID and session ID
    if (launchData?.activity.title === 'Typing Tutor') {
      try {
        // Create a study session and get the session ID
        const groupId = parseInt(selectedGroup)
        const activityId = launchData.activity.id
        const sessionResponse = await createStudySession(groupId, activityId)
        
        // Add both group_id and session_id to the URL
        const typingTutorUrl = new URL(launchData.activity.launch_url)
        typingTutorUrl.searchParams.set('group_id', selectedGroup)
        typingTutorUrl.searchParams.set('session_id', sessionResponse.id.toString())
        
        // Open the typing tutor in a new window
        window.open(typingTutorUrl.toString(), '_blank')
      } catch (error) {
        setError('Failed to create study session')
        console.error(error)
      }
    } else {
      // Handle other activities as before
      window.open(launchData?.activity.launch_url, '_blank')
    }
  }

  if (loading) {
    return <div>Loading...</div>
  }

  if (error) {
    return <div className="text-red-500">{error}</div>
  }

  if (!launchData) {
    return <div>No launch data available</div>
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">{launchData.activity.title}</h1>
      
      <div className="space-y-4">
        <div className="space-y-2">
          <label className="text-sm font-medium">Select Word Group</label>
          <Select onValueChange={setSelectedGroup} value={selectedGroup}>
            <SelectTrigger>
              <SelectValue placeholder="Select a word group" />
            </SelectTrigger>
            <SelectContent>
              {launchData.groups.map((group) => (
                <SelectItem key={group.id} value={group.id.toString()}>
                  {group.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <Button 
          onClick={handleLaunch}
          disabled={!selectedGroup}
          className="w-full"
        >
          Launch Now
        </Button>
      </div>
    </div>
  )
}
