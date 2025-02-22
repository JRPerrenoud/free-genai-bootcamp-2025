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

  const handleLaunch = () => {
    if (!selectedGroup) {
      setError('Please select a word group')
      return
    }

    // For typing tutor, we'll open it in a new window with the group ID
    if (launchData?.activity.title === 'Typing Tutor') {
      const typingTutorUrl = new URL(launchData.activity.launch_url)
      typingTutorUrl.searchParams.set('group_id', selectedGroup)
      window.open(typingTutorUrl.toString(), '_blank')
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
