interface Group {
  id: number
  name: string
  description: string
}

interface Activity {
  id: number
  name: string
  launch_url: string
  created_at: string
}

export interface LastStudySession {
  id: number
  group_id: number
  study_activity_id: number
  created_at: string
  group: Group
  activity: Activity
  review_items_count: number
}

export interface StudyProgress {
  correct_count: number
  wrong_count: number
}

// Backend response type
export interface DashboardStatsResponse {
  total_words: number
  total_groups: number
  total_sessions: number
  total_reviews: number
  correct_reviews: number
  wrong_reviews: number
}

// Frontend display type
export interface QuickStats {
  success_rate: number
  total_sessions: number
  active_groups: number
  study_streak: number
}
