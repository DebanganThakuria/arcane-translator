
import React from 'react';
import { Link } from 'react-router-dom';
import { Novel } from '../types/novel';
import { Badge } from '@/components/ui/badge';
import { formatDistanceToNow } from 'date-fns';

interface NovelCardProps {
  novel: Novel;
  isRecent?: boolean;
}

const NovelCard: React.FC<NovelCardProps> = ({ novel, isRecent }) => {
  return (
    <Link to={`/novel/${novel.id}`} className="group">
      <div className="book-cover aspect-[2/3] transition-all duration-300 group-hover:translate-y-[-0.5rem]">
        <img 
          src={novel.cover || '/placeholder.svg'} 
          alt={novel.title} 
          className="w-full h-full object-cover rounded-md shadow-md group-hover:shadow-lg"
        />
        {isRecent && novel.lastRead && (
          <div className="absolute bottom-0 left-0 right-0 bg-indigo-600/80 text-white text-xs p-1 text-center">
            {Math.round(novel.lastRead.progress * 100)}% read
          </div>
        )}
        {novel.status === 'Ongoing' && (
          <Badge className="absolute top-2 right-2 bg-gradient-to-r from-indigo-600 to-blue-400">
            Ongoing
          </Badge>
        )}
      </div>
      <div className="mt-2">
        <h3 className="font-medium text-sm line-clamp-2 group-hover:text-indigo-600 transition-colors">{novel.title}</h3>
        <p className="text-xs text-muted-foreground mt-1">
          {novel.lastUpdated ? `Updated ${formatDistanceToNow(novel.lastUpdated)} ago` : ''}
        </p>
      </div>
    </Link>
  );
};

export default NovelCard;
