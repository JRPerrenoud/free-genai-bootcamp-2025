export interface StudySession {
  id: number
  activity_name: string
  group_name: string
  start_time: string
  end_time: string
  review_items_count: number
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
