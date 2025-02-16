import { api } from './api'
import { LastStudySession, StudyProgress, DashboardStatsResponse } from '../types/dashboard'

export const dashboardService = {
  getLastStudySession: () => {
    return api.get<LastStudySession>('/api/dashboard/last_study_session')
  },

  getStudyProgress: () => {
    return api.get<StudyProgress>('/api/dashboard/study_progress')
  },

  getQuickStats: () => {
    return api.get<DashboardStatsResponse>('/api/dashboard/quick_stats')
  },
}
