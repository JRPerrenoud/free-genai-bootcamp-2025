export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  current_page: number;
  total_pages: number;
  total_items: number;
  items_per_page: number;
}

export interface LastStudySession {
  id: number;
  group_id: number;
  created_at: string;
  study_activity_id: number;
  group_name: string;
}

export interface StudyProgress {
  total_words: number;
  studied_words: number;
  group_progress: Array<{
    group_id: number;
    group_name: string;
    studied: number;
    total: number;
  }>;
}

export interface QuickStats {
  successRate: number;
  totalSessions: number;
  activeGroups: number;
  studyStreak: number;
}
