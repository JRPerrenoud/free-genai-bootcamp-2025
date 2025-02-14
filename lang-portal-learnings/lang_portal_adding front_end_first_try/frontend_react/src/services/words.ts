import { apiClient } from './api-client';
import type { Word, WordGroup, PaginatedResponse, ApiResponse } from '@/types/words';

export const wordService = {
  getWords: async (page: number = 1) => {
    try {
      const response = await apiClient.get<ApiResponse<PaginatedResponse<Word>>>('/words', {
        params: { page },
      });
      console.log('Raw API response:', JSON.stringify(response.data, null, 2));
      return response.data;
    } catch (error) {
      console.error('Error fetching words:', error);
      return { success: false, data: null };
    }
  },

  getWordById: async (id: number) => {
    try {
      const response = await apiClient.get<ApiResponse<Word>>(`/words/${id}`);
      console.log('Raw API response:', JSON.stringify(response.data, null, 2));
      return response.data;
    } catch (error) {
      console.error('Error fetching word by id:', error);
      return { success: false, data: null };
    }
  },

  createWord: async (word: Omit<Word, 'id'>) => {
    try {
      const response = await apiClient.post<ApiResponse<Word>>('/words', word);
      console.log('Raw create word response:', JSON.stringify(response.data, null, 2));
      return response.data;
    } catch (error) {
      console.error('Error creating word:', error);
      return { success: false, data: null };
    }
  },

  updateWord: async (id: number, word: Omit<Word, 'id'>) => {
    try {
      const response = await apiClient.put<ApiResponse<Word>>(`/words/${id}`, word);
      console.log('Raw update word response:', JSON.stringify(response.data, null, 2));
      return response.data;
    } catch (error) {
      console.error('Error updating word:', error);
      return { success: false, data: null };
    }
  },

  deleteWord: async (id: number) => {
    try {
      const response = await apiClient.delete<ApiResponse<void>>(`/words/${id}`);
      console.log('Raw delete word response:', JSON.stringify(response.data, null, 2));
      return response.data;
    } catch (error) {
      console.error('Error deleting word:', error);
      return { success: false, data: null };
    }
  },
};

export const groupService = {
  getGroups: async (page: number = 1) => {
    const response = await apiClient.get<{ success: boolean; data: PaginatedResponse<WordGroup> }>('/groups', {
      params: { page },
    });
    console.log('Raw API response:', JSON.stringify(response.data, null, 2));
    return response.data;
  },

  getGroupById: async (id: number) => {
    const response = await apiClient.get<{ success: boolean; data: WordGroup }>(`/groups/${id}`);
    console.log('Raw API response:', JSON.stringify(response.data, null, 2));
    return response.data;
  },

  getGroupWords: async (id: number, page: number = 1) => {
    const response = await apiClient.get<{ success: boolean; data: PaginatedResponse<Word> }>(`/groups/${id}/words`, {
      params: { page },
    });
    console.log('Raw API response:', JSON.stringify(response.data, null, 2));
    return response.data;
  },

  createGroup: async (group: Omit<WordGroup, 'id' | 'word_count'>) => {
    const response = await apiClient.post<{ success: boolean; data: WordGroup }>('/groups', group);
    console.log('Raw create group response:', JSON.stringify(response.data, null, 2));
    return response.data;
  },

  updateGroup: async (id: number, group: Omit<WordGroup, 'id' | 'word_count'>) => {
    const response = await apiClient.put<{ success: boolean; data: WordGroup }>(`/groups/${id}`, group);
    console.log('Raw update group response:', JSON.stringify(response.data, null, 2));
    return response.data;
  },

  deleteGroup: async (id: number) => {
    const response = await apiClient.delete<{ success: boolean }>(`/groups/${id}`);
    console.log('Raw delete group response:', JSON.stringify(response.data, null, 2));
    return response.data;
  },

  addWordToGroup: async (groupId: number, wordId: number) => {
    const response = await apiClient.post(`/groups/${groupId}/words`, { word_id: wordId });
    console.log('Raw add word to group response:', JSON.stringify(response.data, null, 2));
    return response.data;
  },

  removeWordFromGroup: async (groupId: number, wordId: number) => {
    const response = await apiClient.delete(`/groups/${groupId}/words/${wordId}`);
    console.log('Raw remove word from group response:', JSON.stringify(response.data, null, 2));
    return response.data;
  },
};
