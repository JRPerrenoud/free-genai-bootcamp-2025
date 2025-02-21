import { api } from './api'
import type { Word } from '../types/api'

interface PaginatedResponse<T> {
  success: boolean
  data: {
    items: T[]
    total_items: number
    current_page: number
    total_pages: number
    items_per_page: number
  }
}

export const wordsService = {
  async getAll(): Promise<Word[]> {
    const response = await api.get<PaginatedResponse<Word>>('/api/words')
    return response.data.items || []
  },

  async getById(id: number): Promise<Word> {
    const response = await api.get<{ success: boolean, data: Word }>(`/api/words/${id}`)
    return response.data
  },

  async create(word: Omit<Word, 'id' | 'created_at' | 'updated_at'>): Promise<Word> {
    const response = await api.post<{ success: boolean, data: Word }>('/api/words', word)
    return response.data
  },

  async update(id: number, word: Partial<Word>): Promise<Word> {
    const response = await api.put<{ success: boolean, data: Word }>(`/api/words/${id}`, word)
    return response.data
  },

  async delete(id: number): Promise<void> {
    await api.delete<{ success: boolean }>(`/api/words/${id}`)
  },
}
