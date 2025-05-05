
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
          className="flex items-center gap-1"
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
            className="flex items-center justify-between hover:bg-accent/50 p-2 rounded-md transition-colors"
          >
            <div className="truncate">
              <span className="text-sm">Ch. {chapter.number}: </span>
              <span className="font-medium">{chapter.title}</span>
              {chapter.wordCount && (
                <span className="text-xs text-muted-foreground ml-2">
                  ({chapter.wordCount} words)
                </span>
              )}
            </div>
            <span className="text-xs text-muted-foreground">
              {formatDistanceToNow(chapter.dateTranslated)} ago
            </span>
          </Link>
        ))}
      </div>
    </div>
  );
};

export default ChapterList;
