import type { StudyActivityResponse, SingleStudyActivityResponse, StudySessionsResponse } from '../types/study-activities'
import { api } from './api'

class StudyActivitiesService {
  async getAll(): Promise<StudyActivityResponse> {
    return api.get<StudyActivityResponse>('/api/study_activities')
  }

  async getById(id: number): Promise<SingleStudyActivityResponse> {
    return api.get<SingleStudyActivityResponse>(`/api/study_activities/${id}`)
  }

  async getSessions(id: number): Promise<StudySessionsResponse> {
    return api.get<StudySessionsResponse>(`/api/study_activities/${id}/sessions`)
  }

  async getActivityDetail(id: number): Promise<SingleStudyActivityResponse> {
    return api.get<SingleStudyActivityResponse>(`/api/study_activities/${id}`)
  }
}

export const studyActivitiesService = new StudyActivitiesService()
