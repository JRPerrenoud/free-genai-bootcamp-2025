import { FC, useState } from 'react'
import { useWords } from '../../hooks/useWords'
import type { Word } from '../../types/api'

const WordsPage: FC = () => {
  const { words, loading, error, addWord } = useWords()
  const [newWord, setNewWord] = useState({
    original: '',
    translated: '',
    pronunciation: '',
    context: '',
    notes: ''
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await addWord(newWord)
      setNewWord({
        original: '',
        translated: '',
        pronunciation: '',
        context: '',
        notes: ''
      })
    } catch (err) {
      console.error('Failed to add word:', err)
    }
  }

  if (loading) return <div className="p-6">Loading...</div>
  if (error) return <div className="p-6 text-red-500">Error: {error.message}</div>

  return (
    <div className="container mx-auto p-6">
      <h1 className="text-3xl font-bold mb-6">Words</h1>
      
      {/* Add Word Form */}
      <form onSubmit={handleSubmit} className="mb-8 space-y-4 max-w-xl">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1">Original</label>
            <input
              type="text"
              value={newWord.original}
              onChange={e => setNewWord(prev => ({ ...prev, original: e.target.value }))}
              className="w-full p-2 border rounded"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Translated</label>
            <input
              type="text"
              value={newWord.translated}
              onChange={e => setNewWord(prev => ({ ...prev, translated: e.target.value }))}
              className="w-full p-2 border rounded"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Pronunciation</label>
            <input
              type="text"
              value={newWord.pronunciation}
              onChange={e => setNewWord(prev => ({ ...prev, pronunciation: e.target.value }))}
              className="w-full p-2 border rounded"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Context</label>
            <input
              type="text"
              value={newWord.context}
              onChange={e => setNewWord(prev => ({ ...prev, context: e.target.value }))}
              className="w-full p-2 border rounded"
            />
          </div>
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Notes</label>
          <textarea
            value={newWord.notes}
            onChange={e => setNewWord(prev => ({ ...prev, notes: e.target.value }))}
            className="w-full p-2 border rounded"
            rows={2}
          />
        </div>
        <button
          type="submit"
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Add Word
        </button>
      </form>

      {/* Words List */}
      <div className="grid gap-4">
        {words.map((word: Word) => (
          <div
            key={word.id}
            className="p-4 border rounded shadow-sm hover:shadow-md transition-shadow"
          >
            <div className="flex justify-between items-start">
              <div>
                <h3 className="font-medium">{word.original}</h3>
                <p className="text-gray-600">{word.translated}</p>
                {word.pronunciation && (
                  <p className="text-sm text-gray-500">/{word.pronunciation}/</p>
                )}
              </div>
              <div className="text-sm text-gray-500">
                {new Date(word.created_at).toLocaleDateString()}
              </div>
            </div>
            {word.context && (
              <p className="mt-2 text-sm text-gray-600">
                <span className="font-medium">Context:</span> {word.context}
              </p>
            )}
            {word.notes && (
              <p className="mt-1 text-sm text-gray-600">
                <span className="font-medium">Notes:</span> {word.notes}
              </p>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}

export default WordsPage
