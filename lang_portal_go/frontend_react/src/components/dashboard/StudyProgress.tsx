import { FC } from 'react'
import type { StudyProgress as StudyProgressType } from '../../types/dashboard'

interface StudyProgressProps {
  data: StudyProgressType | null
  isLoading: boolean
}

const StudyProgress: FC<StudyProgressProps> = ({ data, isLoading }) => {
  if (isLoading) {
    return (
      <div className="p-6 bg-white rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Study Progress</h2>
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-full mb-4"></div>
          <div className="h-8 bg-gray-200 rounded"></div>
        </div>
      </div>
    )
  }

  if (!data) {
    return (
      <div className="p-6 bg-white rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Study Progress</h2>
        <p className="text-gray-600">No progress data available</p>
      </div>
    )
  }

  const totalAnswers = data.correct_count + data.wrong_count
  const accuracy = totalAnswers === 0 ? 0 : (data.correct_count / totalAnswers) * 100

  return (
    <div className="p-6 bg-white rounded-lg shadow">
      <h2 className="text-xl font-semibold mb-4">Study Progress</h2>
      <div className="space-y-4">
        <div>
          <p className="text-gray-600">Total Reviews</p>
          <p className="text-2xl font-semibold">{totalAnswers}</p>
        </div>
        <div>
          <p className="text-gray-600">Accuracy</p>
          <div className="flex items-center space-x-2">
            <div className="flex-1 bg-gray-200 rounded-full h-2.5">
              <div
                className="bg-green-600 h-2.5 rounded-full"
                style={{ width: `${accuracy}%` }}
              ></div>
            </div>
            <span className="text-sm font-medium">{accuracy.toFixed(1)}%</span>
          </div>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <p className="text-gray-600">Correct</p>
            <p className="text-green-600 font-semibold">{data.correct_count}</p>
          </div>
          <div>
            <p className="text-gray-600">Wrong</p>
            <p className="text-red-600 font-semibold">{data.wrong_count}</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default StudyProgress
