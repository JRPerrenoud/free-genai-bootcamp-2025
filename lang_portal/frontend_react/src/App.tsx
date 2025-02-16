import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/common/Layout'
import DashboardPage from './pages/dashboard'
import WordsPage from './pages/words'
import StudyActivitiesPage from './pages/study-activities'
import StudyActivityShowPage from './pages/study-activities/show'

function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          {/* Redirect root to dashboard */}
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          
          {/* Main routes */}
          <Route path="/dashboard" element={<DashboardPage />} />
          
          {/* Study Activities routes */}
          <Route path="/study_activities" element={<StudyActivitiesPage />} />
          <Route path="/study_activities/:id" element={<StudyActivityShowPage />} />
          <Route path="/study_activities/:id/launch" element={<div>Launch Study Activity</div>} />
          
          {/* Words routes */}
          <Route path="/words" element={<WordsPage />} />
          <Route path="/words/:id" element={<div>Word Details</div>} />
          
          {/* Groups routes */}
          <Route path="/groups" element={<div>Groups List</div>} />
          <Route path="/groups/:id" element={<div>Group Details</div>} />
          
          {/* Sessions routes */}
          <Route path="/sessions" element={<div>Sessions List</div>} />
          <Route path="/sessions/:id" element={<div>Session Details</div>} />
          
          {/* 404 route */}
          <Route path="*" element={<div>Page Not Found</div>} />
        </Routes>
      </Layout>
    </BrowserRouter>
  )
}

export default App
