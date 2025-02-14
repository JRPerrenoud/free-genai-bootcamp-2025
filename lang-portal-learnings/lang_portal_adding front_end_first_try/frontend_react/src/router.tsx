import { createBrowserRouter } from 'react-router-dom';
import { RootLayout } from './components/layout/RootLayout';
import DashboardPage from './pages/DashboardPage';
import WordsPage from './pages/WordsPage';
import GroupsPage from './pages/GroupsPage';
import GroupDetailPage from './pages/GroupDetailPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    children: [
      {
        path: '/',
        element: <DashboardPage />,
      },
      {
        path: '/words',
        element: <WordsPage />,
      },
      {
        path: '/groups',
        element: <GroupsPage />,
      },
      {
        path: '/groups/:partOfSpeech/words',
        element: <GroupDetailPage />,
      },
      // Add more routes here as we create them
    ],
  },
]);
