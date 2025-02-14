import { cn } from '../../lib/utils';
import { HTMLAttributes, forwardRef } from 'react';

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  title?: string;
}

const Card = forwardRef<HTMLDivElement, CardProps>(
  ({ className, title, children, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={cn(
          'rounded-lg border bg-white shadow-sm',
          className
        )}
        {...props}
      >
        {title && (
          <div className="border-b px-4 py-3">
            <h3 className="font-semibold">{title}</h3>
          </div>
        )}
        <div className="p-4">
          {children}
        </div>
      </div>
    );
  }
);

Card.displayName = 'Card';

export { Card, type CardProps };
