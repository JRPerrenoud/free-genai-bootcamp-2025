import { cn } from '../../lib/utils';
import { HTMLAttributes } from 'react';
import { Card } from './Card';

interface StatCardProps extends HTMLAttributes<HTMLDivElement> {
  title: string;
  value: number | string;
  icon?: React.ReactNode;
  trend?: number;
}

export function StatCard({ 
  title, 
  value, 
  icon, 
  trend, 
  className,
  ...props 
}: StatCardProps) {
  return (
    <Card
      className={cn('overflow-hidden', className)}
      {...props}
    >
      <div className="flex items-start justify-between">
        <div>
          <p className="text-sm font-medium text-gray-500">{title}</p>
          <h3 className="mt-2 text-2xl font-semibold text-gray-900">{value}</h3>
        </div>
        {icon && (
          <div className="rounded-full bg-gray-100 p-2">
            {icon}
          </div>
        )}
      </div>
      {trend !== undefined && (
        <div className={cn(
          'mt-4 flex items-center text-sm font-medium',
          trend >= 0 ? 'text-green-600' : 'text-red-600'
        )}>
          <span className="mr-1">
            {trend >= 0 ? '↑' : '↓'}
          </span>
          <span>{Math.abs(trend)}%</span>
          <span className="ml-2 text-gray-600">from last period</span>
        </div>
      )}
    </Card>
  );
}
