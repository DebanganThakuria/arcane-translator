
import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Chapter, Novel } from '../types/novel';
import { ArrowLeft, ArrowRight } from 'lucide-react';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { useToast } from '@/components/ui/use-toast';
import { translateChapter } from '../services/translationService';

interface ReaderProps {
  chapter: Chapter;
  novel: Novel;
  hasPreviousChapter: boolean;
  hasNextChapter: boolean;
  onNavigate: (direction: 'next' | 'prev') => void;
}

const Reader: React.FC<ReaderProps> = ({
  chapter,
  novel,
  hasPreviousChapter,
  hasNextChapter,
  onNavigate,
}) => {
  const [sepiaMode, setSepiaMode] = useState(true);
  const [fontSize, setFontSize] = useState(18);
  const { toast } = useToast();
  const navigate = useNavigate();

  useEffect(() => {
    // Scroll to top when chapter changes
    window.scrollTo(0, 0);
    
    // When a chapter is loaded, simulate background translation of the next chapter
    if (hasNextChapter) {
      const nextChapterNumber = chapter.number + 1;
      
      // Simulate background translation - this would be a real API call in production
      const simulateBackgroundTranslation = async () => {
        try {
          console.log(`Background translating chapter ${nextChapterNumber}...`);
          
          // This is where we would actually call the translation service
          // In a real app, we'd store the result in the database
          // translateChapter(novel.id, nextChapterNumber);
        } catch (error) {
          console.error('Error in background translation:', error);
        }
      };
      
      // Wait a few seconds before starting background translation
      const timeoutId = setTimeout(simulateBackgroundTranslation, 5000);
      
      return () => {
        clearTimeout(timeoutId);
      };
    }
  }, [chapter.id, chapter.number, hasNextChapter, novel.id]);

  const handleFontSizeChange = (increase: boolean) => {
    setFontSize(prev => {
      const newSize = increase ? prev + 1 : prev - 1;
      return Math.min(Math.max(newSize, 14), 24); // Limit between 14 and 24
    });
  };

  const handleNavigate = (direction: 'next' | 'prev') => {
    onNavigate(direction);
    
    // Simulate translation notification for the next chapter
    if (direction === 'next' && hasNextChapter) {
      const nextChapterNumber = chapter.number + 1;
      toast({
        title: "Auto-translation",
        description: `Chapter ${nextChapterNumber + 1} is being translated in the background`,
      });
    }
  };

  // Format chapter content to add line breaks after paragraphs
  const formattedContent = chapter.content.replace(/<\/p>/g, '</p><br/>');

  return (
    <div className={`min-h-screen pb-20 ${sepiaMode ? 'sepia-mode' : 'bg-background'}`}>
      <header className="sticky top-0 z-10 backdrop-blur-md bg-background/80 dark:bg-background/80 border-b">
        <div className="container mx-auto p-4 flex justify-between items-center">
          <Button variant="ghost" size="sm" onClick={() => navigate(`/novel/${novel.id}`)}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Novel
          </Button>
          <div className="flex items-center gap-4">
            <div className="flex items-center space-x-2">
              <Switch
                id="sepia-mode"
                checked={sepiaMode}
                onCheckedChange={setSepiaMode}
              />
              <Label htmlFor="sepia-mode">Sepia</Label>
            </div>
            
            <div className="flex items-center space-x-2">
              <Button 
                variant="outline" 
                size="icon" 
                onClick={() => handleFontSizeChange(false)}
                disabled={fontSize <= 14}
              >
                A-
              </Button>
              <Button 
                variant="outline" 
                size="icon" 
                onClick={() => handleFontSizeChange(true)}
                disabled={fontSize >= 24}
              >
                A+
              </Button>
            </div>
          </div>
        </div>
      </header>

      <div className={`reader-container ${sepiaMode ? 'sepia' : 'bg-white'}`}>
        <h1 className="text-2xl font-semibold mb-2">{chapter.title}</h1>
        <p className="text-sm text-muted-foreground mb-6">
          {novel.title} • Chapter {chapter.number}
        </p>
        
        <div 
          className="prose max-w-none"
          style={{ fontSize: `${fontSize}px` }}
          dangerouslySetInnerHTML={{ __html: formattedContent }}
        />
      </div>
      
      <div className="fixed bottom-0 left-0 right-0 p-4 backdrop-blur-md bg-background/80 dark:bg-background/80 border-t">
        <div className="container mx-auto flex justify-between">
          <Button
            variant="outline"
            onClick={() => handleNavigate('prev')}
            disabled={!hasPreviousChapter}
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Previous Chapter
          </Button>
          
          <Button
            onClick={() => handleNavigate('next')}
            disabled={!hasNextChapter}
          >
            Next Chapter
            <ArrowRight className="ml-2 h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
};

export default Reader;
