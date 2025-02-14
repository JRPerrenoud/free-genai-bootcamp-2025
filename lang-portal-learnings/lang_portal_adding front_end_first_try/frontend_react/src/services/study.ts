import { apiClient } from './api-client';
import { ApiResponse, PaginatedResponse } from '../types/api';

interface StudyActivity {
  id: number;
  name: string;
  description: string;
  created_at: string;
}

export const studyService = {
  getActivities: async (page: number = 1) => {
    const response = await apiClient.get<ApiResponse<PaginatedResponse<StudyActivity>>>('/study-activities', {
      params: { page },
    });
    return response.data;
  },
};
