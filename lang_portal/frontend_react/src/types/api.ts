export interface Word {
  id: number
  original: string
  translated: string
  pronunciation?: string
  context?: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface Group {
  id: number
  name: string
  description?: string
  created_at: string
  updated_at: string
  words?: Word[]
}

export interface StudyActivity {
  id: number
  name: string
  description?: string
  activity_type: string
  config: Record<string, any>
  created_at: string
  updated_at: string
}

export interface StudySession {
  id: number
  activity_id: number
  start_time: string
  end_time?: string
  score?: number
  data: Record<string, any>
}
