import type { StudySessionsResponse } from '../types/study-sessions'
import { api } from './api'

class StudySessionsService {
  async getByActivityId(activityId: number): Promise<StudySessionsResponse> {
    return api.get<StudySessionsResponse>(`/api/study_activities/${activityId}/study_sessions`)
  }
}

export const studySessionsService = new StudySessionsService()
