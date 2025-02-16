import type { StudyActivityResponse, SingleStudyActivityResponse } from '../types/study-activities'
import { api } from './api'

class StudyActivitiesService {
  async getAll(): Promise<StudyActivityResponse> {
    return api.get<StudyActivityResponse>('/api/study_activities')
  }

  async getById(id: number): Promise<SingleStudyActivityResponse> {
    return api.get<SingleStudyActivityResponse>(`/api/study_activities/${id}`)
  }
}

export const studyActivitiesService = new StudyActivitiesService()
