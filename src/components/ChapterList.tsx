import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Chapter } from '../types/novel';
import { formatDistanceToNow } from 'date-fns';
import { Button } from '@/components/ui/button';
import { ArrowDown, ArrowUp } from 'lucide-react';

interface ChapterListProps {
  chapters: Chapter[];
  novelId: string;
}

const ChapterList: React.FC<ChapterListProps> = ({ chapters, novelId }) => {
  const [sortAscending, setSortAscending] = useState(true);
  
  if (!chapters.length) {
    return <p className="text-center py-8 text-muted-foreground">No chapters available</p>;
  }

  // Sort chapters by number based on sort direction
  const sortedChapters = [...chapters].sort((a, b) => 
    sortAscending ? a.number - b.number : b.number - a.number
  );

  const toggleSortOrder = () => {
    setSortAscending(!sortAscending);
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
        {sortedChapters.map((chapter) => (
          <Link 
            key={chapter.id} 
            to={`/novel/${novelId}/chapter/${chapter.number}`}
            className="flex items-start hover:bg-indigo-50/50 p-2 rounded-md transition-colors group"
          >
            <div className="flex-1 min-w-0">
              <div className="flex items-baseline">
                <span className="text-sm font-medium text-muted-foreground mr-2 flex-shrink-0">
                  {chapter.number}.
                </span>
                <span className="font-medium truncate">
                  {chapter.title}
                </span>
              </div>
              {chapter.word_count && (
                <div className="text-xs text-muted-foreground mt-1">
                  {chapter.word_count} words
                </div>
              )}
            </div>
            <span className="text-xs text-muted-foreground whitespace-nowrap ml-2 self-center">
              {chapter.date_translated ? formatDistanceToNow(new Date(chapter.date_translated * 1000)) + ' ago' : 'Not translated'}
            </span>
          </Link>
        ))}
      </div>
    </div>
  );
};

export default ChapterList;
