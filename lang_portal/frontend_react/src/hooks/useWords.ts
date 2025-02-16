import { useState, useEffect } from 'react'
import { wordsService } from '../services/words'
import type { Word } from '../types/api'

export function useWords() {
  const [words, setWords] = useState<Word[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    loadWords()
  }, [])

  async function loadWords() {
    try {
      setLoading(true)
      const data = await wordsService.getAll()
      setWords(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Failed to load words'))
    } finally {
      setLoading(false)
    }
  }

  async function addWord(word: Omit<Word, 'id' | 'created_at' | 'updated_at'>) {
    try {
      const newWord = await wordsService.create(word)
      setWords(prev => [...prev, newWord])
      return newWord
    } catch (err) {
      throw err instanceof Error ? err : new Error('Failed to add word')
    }
  }

  async function updateWord(id: number, updates: Partial<Word>) {
    try {
      const updatedWord = await wordsService.update(id, updates)
      setWords(prev => prev.map(word => word.id === id ? updatedWord : word))
      return updatedWord
    } catch (err) {
      throw err instanceof Error ? err : new Error('Failed to update word')
    }
  }

  async function deleteWord(id: number) {
    try {
      await wordsService.delete(id)
      setWords(prev => prev.filter(word => word.id !== id))
    } catch (err) {
      throw err instanceof Error ? err : new Error('Failed to delete word')
    }
  }

  return {
    words,
    loading,
    error,
    addWord,
    updateWord,
    deleteWord,
    refresh: loadWords,
  }
}
