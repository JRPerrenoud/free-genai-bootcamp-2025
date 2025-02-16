import { FC, ReactNode } from 'react'
import { Link, useLocation } from 'react-router-dom'

interface LayoutProps {
  children: ReactNode
}

const Layout: FC<LayoutProps> = ({ children }) => {
  const location = useLocation()

  return (
    <div className="min-h-screen bg-background flex">
      {/* Left Sidebar */}
      <div className="w-64 min-h-screen border-r bg-white">
        <div className="p-6">
          <Link to="/" className="text-xl font-bold block mb-8">
            Language Portal
          </Link>
          <nav className="flex flex-col space-y-4">
            <Link 
              to="/dashboard" 
              className={`hover:text-primary ${location.pathname === '/dashboard' ? 'text-primary font-medium' : ''}`}
            >
              Dashboard
            </Link>
            <Link 
              to="/study_activities" 
              className={`hover:text-primary ${location.pathname === '/study_activities' ? 'text-primary font-medium' : ''}`}
            >
              Study Activities
            </Link>
            <Link 
              to="/words" 
              className={`hover:text-primary ${location.pathname === '/words' ? 'text-primary font-medium' : ''}`}
            >
              Words
            </Link>
            <Link 
              to="/groups" 
              className={`hover:text-primary ${location.pathname === '/groups' ? 'text-primary font-medium' : ''}`}
            >
              Word Groups
            </Link>
            <Link 
              to="/sessions" 
              className={`hover:text-primary ${location.pathname === '/sessions' ? 'text-primary font-medium' : ''}`}
            >
              Sessions
            </Link>
          </nav>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1">
        <main className="h-full">{children}</main>
      </div>
    </div>
  )
}

export default Layout
