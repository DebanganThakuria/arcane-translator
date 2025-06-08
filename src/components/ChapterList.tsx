import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Chapter } from '../types/novel';
import { formatDistanceToNow } from 'date-fns';
import { Button } from '@/components/ui/button';
import { ArrowDown, ArrowUp, BookOpen, Eye } from 'lucide-react';
import { getReadingProgress } from '../utils/readingProgress';
import { Badge } from '@/components/ui/badge';

interface ChapterListProps {
  chapters: Chapter[];
  novelId: string;
}

const ChapterList: React.FC<ChapterListProps> = ({ chapters, novelId }) => {
  const [sortAscending, setSortAscending] = useState(true);
  const [refreshKey, setRefreshKey] = useState(0);
  
  // Force refresh when component mounts to get latest progress
  useEffect(() => {
    setRefreshKey(prev => prev + 1);
  }, [novelId]);
  
  if (!chapters.length) {
    return <p className="text-center py-8 text-muted-foreground">No chapters available</p>;
  }

  // Get reading progress for this novel (force refresh with key)
  const readingProgress = getReadingProgress(novelId);

  // Sort chapters by number based on sort direction
  const sortedChapters = [...chapters].sort((a, b) => 
    sortAscending ? a.number - b.number : b.number - a.number
  );

  const toggleSortOrder = () => {
    setSortAscending(!sortAscending);
  };

  const isLastReadChapter = (chapterNumber: number) => {
    return readingProgress && readingProgress.lastChapter === chapterNumber;
  };

  const getChapterProgressPercentage = (chapterNumber: number) => {
    if (readingProgress && readingProgress.lastChapter === chapterNumber) {
      return readingProgress.progress;
    }
    return 0;
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-sm font-medium text-muted-foreground">
          {chapters.length} Chapters
        </h3>
        <Button 
          variant="outline" 
          size="sm" 
          onClick={toggleSortOrder}
          className="flex items-center gap-1 border-indigo-200 dark:border-indigo-800 text-indigo-600 dark:text-indigo-400 hover:bg-indigo-50 dark:hover:bg-indigo-950/50"
        >
          {sortAscending ? (
            <>
              Ascending <ArrowUp className="h-4 w-4" />
            </>
          ) : (
            <>
              Descending <ArrowDown className="h-4 w-4" />
            </>
          )}
        </Button>
      </div>
      
      <div className="space-y-1">
        {sortedChapters.map((chapter) => {
          const isCurrentlyReading = isLastReadChapter(chapter.number);
          const progressPercentage = getChapterProgressPercentage(chapter.number);
          
          return (
            <Link 
              key={chapter.id} 
              to={`/novel/${novelId}/chapter/${chapter.number}`}
              target="_blank"
              rel="noopener noreferrer"
              className={`flex items-start p-2 rounded-md transition-colors group relative ${
                isCurrentlyReading 
                  ? 'bg-indigo-50 border border-indigo-200 hover:bg-indigo-100' 
                  : 'hover:bg-indigo-50/50'
              }`}
            >
              <div className="flex-1 min-w-0">
                <div className="flex items-baseline">
                  <span className="text-sm font-medium text-muted-foreground mr-2 flex-shrink-0">
                    {chapter.number}.
                  </span>
                  <span className={`font-medium truncate ${isCurrentlyReading ? 'text-indigo-700' : ''}`}>
                    {chapter.title}
                  </span>
                  {isCurrentlyReading && (
                    <div className="ml-2 flex items-center gap-1">
                      <Badge variant="secondary" className="bg-indigo-100 text-indigo-700 text-xs">
                        <BookOpen className="h-3 w-3 mr-1" />
                        Reading
                      </Badge>
                    </div>
                  )}
                </div>
                
                <div className="flex items-center gap-2 mt-1">
                  {isCurrentlyReading && progressPercentage > 0 && (
                    <span className="text-xs text-indigo-600 font-medium">
                      {Math.round(progressPercentage)}% complete
                    </span>
                  )}
                </div>
                
                {/* Progress bar for currently reading chapter */}
                {isCurrentlyReading && progressPercentage > 0 && (
                  <div className="mt-2 w-full bg-indigo-100 rounded-full h-1.5">
                    <div 
                      className="bg-indigo-500 h-1.5 rounded-full transition-all duration-300"
                      style={{ width: `${progressPercentage}%` }}
                    />
                  </div>
                )}
              </div>
              
              <span className="text-xs text-muted-foreground whitespace-nowrap ml-2 self-center">
                {chapter.date_translated ? formatDistanceToNow(new Date(chapter.date_translated * 1000)) + ' ago' : 'Not translated'}
              </span>
            </Link>
          );
        })}
      </div>
    </div>
  );
};

export default ChapterList;
