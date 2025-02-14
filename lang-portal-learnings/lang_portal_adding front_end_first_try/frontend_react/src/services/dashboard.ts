import { apiClient } from './api-client';
import { ApiResponse, LastStudySession, StudyProgress } from '../types/api';

export const dashboardService = {
  getLastStudySession: async () => {
    const response = await apiClient.get<ApiResponse<LastStudySession>>('/dashboard/last_study_session');
    return response.data;
  },

  getStudyProgress: async () => {
    const response = await apiClient.get<ApiResponse<StudyProgress>>('/dashboard/study_progress');
    return response.data;
  },
};
