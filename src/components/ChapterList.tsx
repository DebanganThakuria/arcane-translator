
import React from 'react';
import { Link } from 'react-router-dom';
import { Chapter } from '../types/novel';
import { formatDistanceToNow } from 'date-fns';
import { Card, CardContent } from '@/components/ui/card';

interface ChapterListProps {
  chapters: Chapter[];
  novelId: string;
}

const ChapterList: React.FC<ChapterListProps> = ({ chapters, novelId }) => {
  if (!chapters.length) {
    return <p className="text-center py-8 text-muted-foreground">No chapters available</p>;
  }

  // Sort chapters by number in ascending order
  const sortedChapters = [...chapters].sort((a, b) => a.number - b.number);

  return (
    <div className="space-y-3">
      {sortedChapters.map((chapter) => (
        <Link key={chapter.id} to={`/novel/${novelId}/chapter/${chapter.number}`}>
          <Card className="hover:bg-accent/50 transition-colors">
            <CardContent className="p-4 flex justify-between items-center">
              <div>
                <p className="font-medium">{chapter.title}</p>
                <p className="text-xs text-muted-foreground mt-1">
                  Translated {formatDistanceToNow(chapter.dateTranslated)} ago
                  {chapter.wordCount && ` • ${chapter.wordCount} words`}
                </p>
              </div>
              <div className="text-sm text-muted-foreground">
                Ch. {chapter.number}
              </div>
            </CardContent>
          </Card>
        </Link>
      ))}
    </div>
  );
};

export default ChapterList;
