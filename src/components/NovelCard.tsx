import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Novel } from '../types/novel';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { formatDistanceToNow } from 'date-fns';
import { BookOpen, Clock, Star, TrendingUp, Eye, EyeOff } from 'lucide-react';
import { Progress } from '@/components/ui/progress';

interface NovelCardProps {
  novel: Novel;
  isRecent?: boolean;
  showStats?: boolean;
}

const NovelCard: React.FC<NovelCardProps> = ({ novel, isRecent, showStats = false }) => {
  const [imageLoaded, setImageLoaded] = useState(false);
  const [imageError, setImageError] = useState(false);
  
  const readingProgress = novel.last_read_chapter_number && novel.chapters_count 
    ? Math.round((novel.last_read_chapter_number / novel.chapters_count) * 100)
    : 0;

  const getStatusColor = (status?: string) => {
    switch (status) {
      case 'Ongoing':
        return 'bg-gradient-to-r from-green-500 to-emerald-400';
      case 'Completed':
        return 'bg-gradient-to-r from-blue-500 to-cyan-400';
      default:
        return 'bg-gradient-to-r from-gray-500 to-slate-400';
    }
  };

  const getLanguageFlag = (source: string) => {
    if (source.toLowerCase().includes('chinese') || source.toLowerCase().includes('cn')) return '🇨🇳';
    if (source.toLowerCase().includes('korean') || source.toLowerCase().includes('kr')) return '🇰🇷';
    if (source.toLowerCase().includes('japanese') || source.toLowerCase().includes('jp')) return '🇯🇵';
    return '🌐';
  };

  return (
    <div className="group relative">
      <Link to={`/novel/${novel.id}`} className="block">
        <div className="book-cover aspect-[2/3] relative overflow-hidden">
          {/* Image Container */}
          <div className="relative w-full h-full">
            {!imageError ? (
              <img 
                src={novel.cover || '/placeholder.svg'} 
                alt={novel.title} 
                className={`w-full h-full object-cover ${
                  imageLoaded ? 'opacity-100' : 'opacity-0'
                }`}
                onLoad={() => setImageLoaded(true)}
                onError={() => setImageError(true)}
              />
            ) : (
              <div className="w-full h-full bg-gradient-to-br from-indigo-100 to-purple-100 dark:from-indigo-900 dark:to-purple-900 flex items-center justify-center">
                <BookOpen className="h-12 w-12 text-indigo-400" />
              </div>
            )}
            
            {/* Loading shimmer effect */}
            {!imageLoaded && !imageError && (
              <div className="absolute inset-0 bg-gray-200 dark:bg-gray-700">
                <div className="shimmer absolute inset-0" />
              </div>
            )}
          </div>

          {/* Overlay Content */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100">
            <div className="absolute bottom-4 left-4 right-4">
              {showStats && (
                <div className="flex items-center space-x-2 text-white text-xs mb-2">
                  <div className="flex items-center space-x-1">
                    <Eye className="h-3 w-3" />
                    <span>{novel.chapters_count}</span>
                  </div>
                  <div className="flex items-center space-x-1">
                    <Clock className="h-3 w-3" />
                    <span>{formatDistanceToNow(novel.last_updated * 1000, { addSuffix: true })}</span>
                  </div>
                </div>
              )}
            </div>
          </div>



          {/* Language Flag */}
          <div className="absolute top-2 left-2 text-lg">
            {getLanguageFlag(novel.source)}
          </div>

          {/* Reading Progress for Recent */}
          {isRecent && readingProgress > 0 && (
            <div className="absolute bottom-0 left-0 right-0 bg-black/80 text-white text-xs p-2">
              <div className="flex items-center justify-between mb-1">
                <span>Progress</span>
                <span>{readingProgress}%</span>
              </div>
              <Progress value={readingProgress} className="h-1" />
            </div>
          )}


        </div>

        {/* Novel Info */}
        <div className="mt-3 space-y-2">
          <h3 className="font-medium text-sm line-clamp-2 group-hover:text-indigo-600 dark:group-hover:text-indigo-400">
            {novel.title}
          </h3>
          
          {novel.author && (
            <p className="text-xs text-muted-foreground">
              by {novel.author}
            </p>
          )}

          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <span>
              {novel.chapters_count} chapters
            </span>
            {novel.last_updated && (
              <span>
                {formatDistanceToNow(novel.last_updated * 1000, { addSuffix: true })}
              </span>
            )}
          </div>

          {/* Genres */}
          {novel.genres && novel.genres.length > 0 && (
            <div className="flex flex-wrap gap-1">
              {novel.genres.slice(0, 2).map((genre, index) => (
                <Badge 
                  key={index} 
                  variant="secondary" 
                  className="text-xs px-2 py-0.5 bg-indigo-50 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-300"
                >
                  {genre}
                </Badge>
              ))}
              {novel.genres.length > 2 && (
                <Badge variant="secondary" className="text-xs px-2 py-0.5">
                  +{novel.genres.length - 2}
                </Badge>
              )}
            </div>
          )}
        </div>
      </Link>

      {/* Quick Actions (visible on hover) */}
              <div className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 space-y-2">
        {isRecent && novel.last_read_chapter_number && (
          <Link to={`/novel/${novel.id}/chapter/${novel.last_read_chapter_number}`}>
            <Button size="sm" className="btn-primary text-xs">
              Continue
            </Button>
          </Link>
        )}
      </div>
    </div>
  );
};

export default NovelCard;
