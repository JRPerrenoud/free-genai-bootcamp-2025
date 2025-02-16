import { FC } from 'react'
import type { DashboardStatsResponse, QuickStats as QuickStatsDisplay } from '../../types/dashboard'

interface QuickStatsProps {
  data: DashboardStatsResponse | null
  isLoading: boolean
}

const StatCard: FC<{ label: string; value: number | string; icon?: string }> = ({ 
  label, 
  value,
  icon 
}) => (
  <div className="flex items-center p-4 bg-white rounded-lg shadow">
    {icon && (
      <div className="mr-4 text-blue-600">
        <i className={`fas ${icon} text-xl`}></i>
      </div>
    )}
    <div>
      <p className="text-gray-600 text-sm">{label}</p>
      <p className="text-2xl font-semibold mt-1">{value}</p>
    </div>
  </div>
)

const QuickStats: FC<QuickStatsProps> = ({ data, isLoading }) => {
  if (isLoading) {
    return (
      <div className="p-6 bg-white rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Quick Stats</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
              <div className="h-8 bg-gray-200 rounded w-3/4"></div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  if (!data) {
    return (
      <div className="p-6 bg-white rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Quick Stats</h2>
        <p className="text-gray-600">No statistics available</p>
      </div>
    )
  }

  // Transform backend data into display format
  const displayStats: QuickStatsDisplay = {
    // Calculate success rate from correct and total reviews
    success_rate: data.total_reviews === 0 
      ? 0 
      : (data.correct_reviews / data.total_reviews) * 100,
    
    // Total sessions directly from backend
    total_sessions: data.total_sessions,
    
    // Active groups is total_groups since these are the ones with study sessions
    active_groups: data.total_groups,
    
    // For study streak, we'll use total_sessions as a proxy for now
    // In a real app, we'd want to track this properly in the backend
    study_streak: Math.min(data.total_sessions, 7) // Cap at 7 for demo
  }

  return (
    <div className="p-6 bg-white rounded-lg shadow">
      <h2 className="text-xl font-semibold mb-4">Quick Stats</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <StatCard
          label="Success Rate"
          value={`${displayStats.success_rate.toFixed(1)}%`}
          icon="fa-chart-line"
        />
        <StatCard
          label="Study Sessions"
          value={displayStats.total_sessions}
          icon="fa-book"
        />
        <StatCard
          label="Active Groups"
          value={displayStats.active_groups}
          icon="fa-users"
        />
        <StatCard
          label="Study Streak"
          value={`${displayStats.study_streak} days`}
          icon="fa-fire"
        />
      </div>
    </div>
  )
}

export default QuickStats
