import type { Word, WordGroup } from '@/types/words';

export const mockWords: Word[] = [
  {
    id: 1,
    text: 'casa',
    translation: 'house',
    pronunciation: 'kah-sah',
    notes: 'Common noun for house or home',
    created_at: '2025-01-15T10:00:00Z',
    group_ids: [1],
  },
  {
    id: 2,
    text: 'perro',
    translation: 'dog',
    pronunciation: 'peh-rro',
    notes: 'Common pet animal',
    created_at: '2025-01-16T11:00:00Z',
    group_ids: [1, 2],
  },
  {
    id: 3,
    text: 'libro',
    translation: 'book',
    pronunciation: 'lee-bro',
    notes: 'Reading material',
    created_at: '2025-01-17T12:00:00Z',
    group_ids: [2],
  },
];

export const mockGroups: WordGroup[] = [
  {
    id: 1,
    name: 'Basic Nouns',
    description: 'Common everyday nouns',
    word_count: 10,
    created_at: '2025-01-10T09:00:00Z',
  },
  {
    id: 2,
    name: 'Animals',
    description: 'Common animal names',
    word_count: 5,
    created_at: '2025-01-11T10:00:00Z',
  },
  {
    id: 3,
    name: 'Verbs',
    description: 'Common verbs',
    word_count: 15,
    created_at: '2025-01-12T11:00:00Z',
  },
];
