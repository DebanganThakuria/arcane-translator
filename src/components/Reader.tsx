
import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Chapter, Novel } from '../types/novel';
import { ArrowLeft, ArrowRight, Loader2 } from 'lucide-react';
import { translateChapter } from '../services/translationService';
import { useToast } from '@/components/ui/use-toast';

interface ReaderProps {
  chapter: Chapter;
  novel: Novel;
  hasPreviousChapter: boolean;
  hasNextChapter: boolean;
  onNavigate: (direction: 'next' | 'prev') => void;
  onTranslateNext?: () => void;
}

const Reader: React.FC<ReaderProps> = ({
  chapter,
  novel,
  hasPreviousChapter,
  hasNextChapter,
  onNavigate,
  onTranslateNext,
}) => {
  const [theme, setTheme] = useState<'light' | 'sepia' | 'dark'>('sepia');
  const [fontSize, setFontSize] = useState(18);
  const [isTranslating, setIsTranslating] = useState(false);
  const navigate = useNavigate();
  const { toast } = useToast();

  useEffect(() => {
    // Scroll to top when chapter changes
    window.scrollTo(0, 0);
  }, [chapter.id]);

  const handleFontSizeChange = (increase: boolean) => {
    setFontSize(prev => {
      const newSize = increase ? prev + 1 : prev - 1;
      return Math.min(Math.max(newSize, 14), 24); // Limit between 14 and 24
    });
  };

  const handleNavigate = (direction: 'next' | 'prev') => {
    onNavigate(direction);
  };

  const handleTranslateNext = async () => {
    setIsTranslating(true);
    try {
      await translateChapter(novel.id);
      toast({
        title: "Success",
        description: "Next chapter has been translated successfully!",
      });
      if (onTranslateNext) {
        onTranslateNext();
      }
    } catch (error) {
      console.error('Error translating chapter:', error);
      toast({
        title: "Error",
        description: error instanceof Error ? error.message : "Failed to translate chapter",
        variant: "destructive",
      });
    } finally {
      setIsTranslating(false);
    }
  };

  // Format chapter content to add line breaks between paragraphs
  const formattedContent = chapter.content.replace(/<\/p>/g, '</p><br/>');

  // Apply theme classes to the root element
  const themeClasses = {
    light: 'bg-white text-foreground',
    sepia: 'sepia-mode',
    dark: 'bg-gray-900 text-gray-100',
  };

  return (
    <div className={`min-h-screen pb-20 ${themeClasses[theme]}`}>
      <header className="sticky top-0 z-10 backdrop-blur-md bg-background/80 dark:bg-background/80 border-b">
        <div className="container mx-auto p-4 flex justify-between items-center">
          <Button variant="ghost" size="sm" onClick={() => navigate(`/novel/${novel.id}`)}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Novel
          </Button>
          <div className="flex items-center gap-4">
            <div className="flex items-center space-x-2">
              <div className="flex items-center space-x-2 rounded-full bg-secondary p-1">
                <button
                  onClick={() => setTheme('light')}
                  className={`h-6 w-6 rounded-full ${theme === 'light' ? 'bg-primary' : 'bg-transparent'}`}
                  aria-label="Light theme"
                  style={{ backgroundColor: theme === 'light' ? 'hsl(222.2 47.4% 11.2%)' : 'transparent' }}
                />
                <button
                  onClick={() => setTheme('sepia')}
                  className={`h-6 w-6 rounded-full ${theme === 'sepia' ? 'bg-primary' : 'bg-transparent'}`}
                  aria-label="Sepia theme"
                  style={{ backgroundColor: theme === 'sepia' ? 'hsl(35 60% 60%)' : 'transparent' }}
                />
                <button
                  onClick={() => setTheme('dark')}
                  className={`h-6 w-6 rounded-full ${theme === 'dark' ? 'bg-primary' : 'bg-transparent'}`}
                  aria-label="Dark theme"
                  style={{ backgroundColor: theme === 'dark' ? 'hsl(222.2 84% 4.9%)' : 'transparent' }}
                />
              </div>
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

      <div className={`reader-container ${theme === 'sepia' ? 'sepia' : theme === 'dark' ? 'bg-gray-900' : 'bg-white'}`}>
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
          
          {hasNextChapter ? (
            <Button
              onClick={() => handleNavigate('next')}
            >
              Next Chapter
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          ) : (
            <Button
              onClick={handleTranslateNext}
              disabled={isTranslating}
              className="bg-green-600 hover:bg-green-700"
            >
              {isTranslating ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Translating...
                </>
              ) : (
                <>
                  Translate Next Chapter
                  <ArrowRight className="ml-2 h-4 w-4" />
                </>
              )}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
};

export default Reader;
