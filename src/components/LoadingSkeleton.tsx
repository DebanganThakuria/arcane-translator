import React from 'react';

interface LoadingSkeletonProps {
  variant?: 'card' | 'list' | 'text' | 'avatar' | 'chapter';
  count?: number;
  className?: string;
}

const LoadingSkeleton: React.FC<LoadingSkeletonProps> = ({ 
  variant = 'card', 
  count = 1, 
  className = '' 
}) => {
  const skeletons = Array.from({ length: count }, (_, index) => {
    switch (variant) {
      case 'card':
        return (
          <div key={index} className={`glass-card rounded-lg p-4 space-y-3 ${className}`}>
            <div className="skeleton h-48 w-full rounded-md" />
            <div className="space-y-2">
              <div className="skeleton h-4 w-3/4 rounded" />
              <div className="skeleton h-3 w-1/2 rounded" />
            </div>
          </div>
        );
      
      case 'list':
        return (
          <div key={index} className={`flex items-center space-x-4 p-4 ${className}`}>
            <div className="skeleton h-16 w-16 rounded-md" />
            <div className="space-y-2 flex-1">
              <div className="skeleton h-4 w-3/4 rounded" />
              <div className="skeleton h-3 w-1/2 rounded" />
            </div>
          </div>
        );
      
      case 'text':
        return (
          <div key={index} className={`space-y-2 ${className}`}>
            <div className="skeleton h-4 w-full rounded" />
            <div className="skeleton h-4 w-5/6 rounded" />
            <div className="skeleton h-4 w-4/6 rounded" />
          </div>
        );
      
      case 'avatar':
        return (
          <div key={index} className={`flex items-center space-x-3 ${className}`}>
            <div className="skeleton h-10 w-10 rounded-full" />
            <div className="space-y-1">
              <div className="skeleton h-3 w-24 rounded" />
              <div className="skeleton h-2 w-16 rounded" />
            </div>
          </div>
        );
      
      case 'chapter':
        return (
          <div key={index} className={`chapter-card space-y-2 ${className}`}>
            <div className="flex justify-between items-center">
              <div className="skeleton h-5 w-32 rounded" />
              <div className="skeleton h-4 w-16 rounded" />
            </div>
            <div className="skeleton h-3 w-full rounded" />
            <div className="skeleton h-3 w-3/4 rounded" />
          </div>
        );
      
      default:
        return (
          <div key={index} className={`skeleton h-20 w-full rounded ${className}`} />
        );
    }
  });

  return <>{skeletons}</>;
};

export default LoadingSkeleton; 