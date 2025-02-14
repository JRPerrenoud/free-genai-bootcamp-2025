import { useState, useEffect } from 'react';
import type { Word } from '@/types/words';

const PARTS_OF_SPEECH = [
  'noun',
  'verb',
  'adjective',
  'adverb',
  'pronoun',
  'preposition',
  'conjunction',
  'interjection',
  'article',
  'number',
  'phrase'
] as const;

interface WordFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (word: Omit<Word, 'id'>) => Promise<boolean>;
  initialData?: Word;
}

export function WordFormModal({ isOpen, onClose, onSubmit, initialData }: WordFormModalProps) {
  const [formData, setFormData] = useState<Omit<Word, 'id'>>({
    spanish: '',
    english: '',
    part_of_speech: PARTS_OF_SPEECH[0],
  });

  useEffect(() => {
    if (initialData) {
      setFormData({
        spanish: initialData.spanish,
        english: initialData.english,
        part_of_speech: initialData.part_of_speech,
      });
    } else {
      setFormData({
        spanish: '',
        english: '',
        part_of_speech: PARTS_OF_SPEECH[0],
      });
    }
  }, [initialData, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const success = await onSubmit(formData);
      if (success) {
        onClose();
        setFormData({
          spanish: '',
          english: '',
          part_of_speech: PARTS_OF_SPEECH[0],
        });
      }
    } catch (error) {
      console.error('Error submitting form:', error);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-full max-w-md">
        <h2 className="text-xl font-semibold mb-4">
          {initialData ? 'Edit Word' : 'Add New Word'}
        </h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="spanish" className="block text-sm font-medium text-gray-700">
              Spanish
            </label>
            <input
              type="text"
              id="spanish"
              value={formData.spanish}
              onChange={(e) => setFormData({ ...formData, spanish: e.target.value })}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              required
            />
          </div>

          <div>
            <label htmlFor="english" className="block text-sm font-medium text-gray-700">
              English
            </label>
            <input
              type="text"
              id="english"
              value={formData.english}
              onChange={(e) => setFormData({ ...formData, english: e.target.value })}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              required
            />
          </div>

          <div>
            <label htmlFor="part_of_speech" className="block text-sm font-medium text-gray-700">
              Part of Speech
            </label>
            <select
              id="part_of_speech"
              value={formData.part_of_speech}
              onChange={(e) => setFormData({ ...formData, part_of_speech: e.target.value })}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              required
            >
              {PARTS_OF_SPEECH.map((pos) => (
                <option key={pos} value={pos}>
                  {pos.charAt(0).toUpperCase() + pos.slice(1)}
                </option>
              ))}
            </select>
          </div>

          <div className="flex justify-end space-x-3 mt-6">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 text-sm font-medium text-white bg-blue-500 hover:bg-blue-600 rounded-md"
            >
              {initialData ? 'Save Changes' : 'Add Word'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
