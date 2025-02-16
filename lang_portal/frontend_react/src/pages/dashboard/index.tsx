import { FC, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import LastStudySession from '../../components/dashboard/LastStudySession'
import StudyProgress from '../../components/dashboard/StudyProgress'
import QuickStats from '../../components/dashboard/QuickStats'
import { dashboardService } from '../../services/dashboard'
import type { 
  LastStudySession as LastStudySessionType, 
  StudyProgress as StudyProgressType, 
  DashboardStatsResponse
} from '../../types/dashboard'

const Dashboard: FC = () => {
  const [lastSession, setLastSession] = useState<LastStudySessionType | null>(null)
  const [studyProgress, setStudyProgress] = useState<StudyProgressType | null>(null)
  const [quickStats, setQuickStats] = useState<DashboardStatsResponse | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        const [lastSessionRes, studyProgressRes, quickStatsRes] = await Promise.all([
          dashboardService.getLastStudySession(),
          dashboardService.getStudyProgress(),
          dashboardService.getQuickStats()
        ])

        setLastSession(lastSessionRes.data)
        setStudyProgress(studyProgressRes.data)
        setQuickStats(quickStatsRes.data)
      } catch (error) {
        console.error('Error fetching dashboard data:', error)
      } finally {
        setIsLoading(false)
      }
    }

    fetchDashboardData()
  }, [])

  return (
    <div className="p-8">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <LastStudySession data={lastSession} isLoading={isLoading} />
        <StudyProgress data={studyProgress} isLoading={isLoading} />
      </div>
      <div className="mt-6">
        <QuickStats data={quickStats} isLoading={isLoading} />
      </div>
      <div className="mt-8 text-center">
        <Link
          to="/study_activities"
          className="inline-block bg-blue-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-blue-700 transition-colors"
        >
          Start Studying
        </Link>
      </div>
    </div>
  )
}

export default Dashboard
