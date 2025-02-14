import { apiClient } from './api-client';
import { ApiResponse, GroupDetailResponse, GroupsResponse, GroupWordsResponse } from '../types/words';

export const groupService = {
  getGroups: async () => {
    const response = await apiClient.get<GroupsResponse>('/groups');
    return response.data;
  },

  getGroupByPartOfSpeech: async (partOfSpeech: string) => {
    const response = await apiClient.get<GroupDetailResponse>(`/groups/${partOfSpeech}`);
    return response.data;
  },

  getGroupWords: async (partOfSpeech: string, page: number = 1) => {
    const response = await apiClient.get<GroupWordsResponse>(`/groups/${partOfSpeech}/words`, {
      params: { page },
    });
    return response.data;
  },
};
