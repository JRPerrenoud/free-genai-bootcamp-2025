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
