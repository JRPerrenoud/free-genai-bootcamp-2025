import { api } from './api'
import type { StudyActivityResponse, SingleStudyActivityResponse } from '../types/study-activities'

export const studyActivitiesService = {
  getAll: () => {
    return api.get<StudyActivityResponse>('/api/study_activities')
  },

  getById: (id: number) => {
    return api.get<SingleStudyActivityResponse>(`/api/study_activities/${id}`)
  },
}
