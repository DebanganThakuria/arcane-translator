import React from 'react';
import NovelCard from './NovelCard';
import { Novel } from '../types/novel';

interface NovelGridProps {
  novels: Novel[];
  recent?: boolean;
}

const NovelGrid: React.FC<NovelGridProps> = ({ novels, recent = false }) => {
  if (!novels.length) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground">No novels found</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-6">
      {novels.map(novel => (
        <NovelCard key={novel.id} novel={novel} isRecent={recent} />
      ))}
    </div>
  );
};

export default NovelGrid;
