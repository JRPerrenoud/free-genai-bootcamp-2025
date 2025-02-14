import { Link } from 'react-router-dom';

export function Navbar() {
  return (
    <nav className="bg-white shadow-sm h-16">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-full">
        <div className="flex justify-between items-center h-full">
          <div className="flex items-center">
            <span className="text-xl font-semibold text-gray-900">
              Language Learning Portal
            </span>
          </div>
        </div>
      </div>
    </nav>
  );
}
