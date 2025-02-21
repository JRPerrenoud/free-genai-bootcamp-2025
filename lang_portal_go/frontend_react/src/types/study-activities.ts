export interface StudyActivity {
  id: number
  name: string
  thumbnail_url: string | null
  launch_url: string
  description: string | null
  created_at: string
}

export interface StudyActivityResponse {
  success: boolean
  data: {
    items: StudyActivity[]
    total_items: number
    current_page: number
    total_pages: number
    items_per_page: number
  }
  error?: string
}

export interface SingleStudyActivityResponse {
  success: boolean
  data: StudyActivity
  error?: string
}

export interface StudySession {
  id: number
  group_id: number
  study_activity_id: number
  created_at: string
  group: {
    id: number
    name: string
    description: string
  }
  review_items_count: number
}

export interface StudyActivityDetail {
  activity_id: number
  activity_name: string
  thumbnail_url: string
  description: string
  launch_button: {
    text: string
    url: string
  }
  study_sessions: StudySession[]
}

export interface StudyActivityDetailResponse {
  success: boolean
  data: StudyActivityDetail
  error?: string
}

export interface StudySessionsResponse {
  success: boolean
  data: {
    total_items: number
    current_page: number
    total_pages: number
    items_per_page: number
    items: StudySession[]
  }
  error?: string
}
