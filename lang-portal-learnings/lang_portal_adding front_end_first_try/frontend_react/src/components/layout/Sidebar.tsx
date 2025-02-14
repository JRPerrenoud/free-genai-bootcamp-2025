import { Link } from 'react-router-dom';

export function Sidebar() {
  return (
    <aside className="w-64 bg-white border-r">
      <nav className="p-4">
        <ul className="space-y-2">
          <li>
            <Link
              to="/"
              className="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md"
            >
              Dashboard
            </Link>
          </li>
          <li>
            <Link
              to="/words"
              className="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md"
            >
              Words
            </Link>
          </li>
          <li>
            <Link
              to="/groups"
              className="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md"
            >
              Groups
            </Link>
          </li>
        </ul>
      </nav>
    </aside>
  );
}
