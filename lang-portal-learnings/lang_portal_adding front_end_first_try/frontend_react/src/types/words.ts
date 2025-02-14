export interface Word {
  id: number;
  spanish: string;
  english: string;
  part_of_speech: string;
  correct_count?: number;
  wrong_count?: number;
}

export interface WordGroup {
  name: string;
  word_count: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  current_page: number;
  total_pages: number;
  total_items: number;
  items_per_page: number;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
}

export interface GroupDetailResponse extends ApiResponse<WordGroup> {}

export interface GroupWordsResponse extends ApiResponse<PaginatedResponse<Word>> {}

export type WordsResponse = ApiResponse<PaginatedResponse<Word>>;
export type GroupsResponse = ApiResponse<{items: WordGroup[]}>;
