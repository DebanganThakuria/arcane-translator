import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Chapter, Novel } from '../types/novel';
import { ArrowLeft, ArrowRight, Loader2, Settings, BookOpen, Clock, Eye, Bookmark, Share2 } from 'lucide-react';
import { translateChapter } from '../services/translationService';
import { useToast } from '@/components/ui/use-toast';
import { Progress } from '@/components/ui/progress';
import { Slider } from '@/components/ui/slider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Badge } from '@/components/ui/badge';

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
  const [lineHeight, setLineHeight] = useState(1.6);
  const [isTranslating, setIsTranslating] = useState(false);
  const [readingProgress, setReadingProgress] = useState(0);
  const [readingTime, setReadingTime] = useState(0);
  const [isBookmarked, setIsBookmarked] = useState(false);
  const [showSettings, setShowSettings] = useState(false);
  const navigate = useNavigate();
  const { toast } = useToast();
  const contentRef = useRef<HTMLDivElement>(null);
  const startTimeRef = useRef<number>(Date.now());

  useEffect(() => {
    // Scroll to top when chapter changes
    window.scrollTo(0, 0);
    startTimeRef.current = Date.now();
    setReadingProgress(0);
    setReadingTime(0);

    // Load saved preferences
    const savedTheme = localStorage.getItem('readerTheme') as 'light' | 'sepia' | 'dark';
    const savedFontSize = localStorage.getItem('readerFontSize');
    const savedLineHeight = localStorage.getItem('readerLineHeight');
    
    if (savedTheme) setTheme(savedTheme);
    if (savedFontSize) setFontSize(parseInt(savedFontSize));
    if (savedLineHeight) setLineHeight(parseFloat(savedLineHeight));

    // Check if chapter is bookmarked
    const bookmarks = JSON.parse(localStorage.getItem('bookmarks') || '[]');
    setIsBookmarked(bookmarks.includes(`${novel.id}-${chapter.number}`));
  }, [chapter.id, novel.id, chapter.number]);

  // Track reading progress
  useEffect(() => {
    const handleScroll = () => {
      if (contentRef.current) {
        const { scrollTop, scrollHeight, clientHeight } = document.documentElement;
        const progress = Math.min((scrollTop / (scrollHeight - clientHeight)) * 100, 100);
        setReadingProgress(progress);
      }
    };

    const timer = setInterval(() => {
      setReadingTime(Math.floor((Date.now() - startTimeRef.current) / 1000));
    }, 1000);

    window.addEventListener('scroll', handleScroll);
    return () => {
      window.removeEventListener('scroll', handleScroll);
      clearInterval(timer);
    };
  }, []);

  const handleFontSizeChange = (value: number[]) => {
    const newSize = value[0];
    setFontSize(newSize);
    localStorage.setItem('readerFontSize', newSize.toString());
  };

  const handleLineHeightChange = (value: number[]) => {
    const newHeight = value[0];
    setLineHeight(newHeight);
    localStorage.setItem('readerLineHeight', newHeight.toString());
  };

  const handleThemeChange = (newTheme: 'light' | 'sepia' | 'dark') => {
    setTheme(newTheme);
    localStorage.setItem('readerTheme', newTheme);
  };

  const handleNavigate = (direction: 'next' | 'prev') => {
    // Save reading progress
    const progress = {
      novelId: novel.id,
      chapterNumber: chapter.number,
      readingTime: Math.floor((Date.now() - startTimeRef.current) / 1000),
      progress: readingProgress,
      timestamp: Date.now()
    };
    
    const savedProgress = JSON.parse(localStorage.getItem('readingProgress') || '[]');
    const existingIndex = savedProgress.findIndex((p: any) => 
      p.novelId === novel.id && p.chapterNumber === chapter.number
    );
    
    if (existingIndex >= 0) {
      savedProgress[existingIndex] = progress;
    } else {
      savedProgress.push(progress);
    }
    
    localStorage.setItem('readingProgress', JSON.stringify(savedProgress.slice(-100))); // Keep last 100 entries
    
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

  const toggleBookmark = () => {
    const bookmarkId = `${novel.id}-${chapter.number}`;
    const bookmarks = JSON.parse(localStorage.getItem('bookmarks') || '[]');
    
    if (isBookmarked) {
      const filtered = bookmarks.filter((id: string) => id !== bookmarkId);
      localStorage.setItem('bookmarks', JSON.stringify(filtered));
      setIsBookmarked(false);
      toast({
        title: "Bookmark removed",
        description: "Chapter removed from bookmarks",
      });
    } else {
      bookmarks.push(bookmarkId);
      localStorage.setItem('bookmarks', JSON.stringify(bookmarks));
      setIsBookmarked(true);
      toast({
        title: "Bookmarked",
        description: "Chapter added to bookmarks",
      });
    }
  };

  const shareChapter = async () => {
    const url = window.location.href;
    const text = `Reading "${chapter.title}" from ${novel.title}`;
    
    if (navigator.share) {
      try {
        await navigator.share({ title: text, url });
      } catch (error) {
        // Fallback to clipboard
        navigator.clipboard.writeText(url);
        toast({
          title: "Link copied",
          description: "Chapter link copied to clipboard",
        });
      }
    } else {
      navigator.clipboard.writeText(url);
      toast({
        title: "Link copied",
        description: "Chapter link copied to clipboard",
      });
    }
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const estimatedReadingTime = Math.ceil((chapter.word_count || 1000) / 200); // 200 WPM average

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
      {/* Reading Progress Bar */}
      <div className="fixed top-0 left-0 right-0 z-50">
        <Progress value={readingProgress} className="h-1 rounded-none" />
      </div>

      <header className="sticky top-0 z-10 backdrop-blur-md bg-background/80 dark:bg-background/80 border-b">
        <div className="container mx-auto p-4">
          <div className="flex justify-between items-center">
            <div className="flex items-center gap-4">
              <Button variant="ghost" size="sm" onClick={() => navigate(`/novel/${novel.id}`)}>
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Novel
              </Button>
              
              {/* Chapter Info */}
              <div className="hidden sm:flex items-center space-x-4 text-sm text-muted-foreground">
                <div className="flex items-center space-x-1">
                  <BookOpen className="h-4 w-4" />
                  <span>Chapter {chapter.number}</span>
                </div>
                <div className="flex items-center space-x-1">
                  <Clock className="h-4 w-4" />
                  <span>{formatTime(readingTime)} / ~{estimatedReadingTime}m</span>
                </div>
                <div className="flex items-center space-x-1">
                  <Eye className="h-4 w-4" />
                  <span>{Math.round(readingProgress)}%</span>
                </div>
              </div>
            </div>

            <div className="flex items-center gap-2">
              {/* Quick Actions */}
              <Button
                variant="ghost"
                size="icon"
                onClick={toggleBookmark}
                className={isBookmarked ? 'text-yellow-500' : ''}
              >
                <Bookmark className={`h-4 w-4 ${isBookmarked ? 'fill-current' : ''}`} />
              </Button>

              <Button variant="ghost" size="icon" onClick={shareChapter}>
                <Share2 className="h-4 w-4" />
              </Button>

              {/* Reader Settings */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon">
                    <Settings className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-80">
                  <DropdownMenuLabel>Reading Settings</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  
                  <div className="p-4 space-y-4">
                    {/* Theme Selection */}
                    <div>
                      <label className="text-sm font-medium mb-2 block">Theme</label>
                      <div className="flex items-center space-x-2">
                        <button
                          onClick={() => handleThemeChange('light')}
                          className={`h-8 w-8 rounded-full border-2 ${theme === 'light' ? 'border-primary' : 'border-gray-300'}`}
                          style={{ backgroundColor: 'white' }}
                          aria-label="Light theme"
                        />
                        <button
                          onClick={() => handleThemeChange('sepia')}
                          className={`h-8 w-8 rounded-full border-2 ${theme === 'sepia' ? 'border-primary' : 'border-gray-300'}`}
                          style={{ backgroundColor: '#FEF7CD' }}
                          aria-label="Sepia theme"
                        />
                        <button
                          onClick={() => handleThemeChange('dark')}
                          className={`h-8 w-8 rounded-full border-2 ${theme === 'dark' ? 'border-primary' : 'border-gray-300'}`}
                          style={{ backgroundColor: '#1f2937' }}
                          aria-label="Dark theme"
                        />
                      </div>
                    </div>

                    {/* Font Size */}
                    <div>
                      <label className="text-sm font-medium mb-2 block">
                        Font Size: {fontSize}px
                      </label>
                      <Slider
                        value={[fontSize]}
                        onValueChange={handleFontSizeChange}
                        max={28}
                        min={12}
                        step={1}
                        className="w-full"
                      />
                    </div>

                    {/* Line Height */}
                    <div>
                      <label className="text-sm font-medium mb-2 block">
                        Line Height: {lineHeight}
                      </label>
                      <Slider
                        value={[lineHeight]}
                        onValueChange={handleLineHeightChange}
                        max={2.5}
                        min={1.2}
                        step={0.1}
                        className="w-full"
                      />
                    </div>
                  </div>
                </DropdownMenuContent>
              </DropdownMenu>

              {!hasNextChapter && (
                <Button
                  onClick={handleTranslateNext}
                  disabled={isTranslating}
                  className="btn-primary"
                >
                  {isTranslating ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Translating...
                    </>
                  ) : (
                    <>
                      Translate Next
                      <ArrowRight className="ml-2 h-4 w-4" />
                    </>
                  )}
                </Button>
              )}
            </div>
          </div>

          {/* Mobile Chapter Info */}
          <div className="sm:hidden mt-4 flex items-center justify-between text-sm text-muted-foreground">
            <div className="flex items-center space-x-4">
              <span>Ch. {chapter.number}</span>
              <span>{formatTime(readingTime)}</span>
              <span>{Math.round(readingProgress)}%</span>
            </div>
            <Badge variant="secondary">
              {novel.title}
            </Badge>
          </div>
        </div>
      </header>

      <div 
        ref={contentRef}
        className={`reader-container ${theme === 'sepia' ? 'sepia' : theme === 'dark' ? 'dark' : 'bg-white'}`}
      >
        <div className="mb-8">
          <h1 className="text-2xl sm:text-3xl font-bold mb-4">{chapter.title}</h1>
          <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
            <span>{novel.title}</span>
            <span>•</span>
            <span>Chapter {chapter.number}</span>
            {chapter.word_count && (
              <>
                <span>•</span>
                <span>{chapter.word_count.toLocaleString()} words</span>
              </>
            )}
            <span>•</span>
            <span>~{estimatedReadingTime} min read</span>
          </div>
        </div>
        
        <div 
          className="prose prose-lg max-w-none dark:prose-invert"
          style={{ 
            fontSize: `${fontSize}px`,
            lineHeight: lineHeight
          }}
          dangerouslySetInnerHTML={{ __html: formattedContent }}
        />

        {/* Chapter End Actions */}
        <div className="mt-12 pt-8 border-t border-gray-200 dark:border-gray-700">
          <div className="text-center space-y-4">
            <p className="text-muted-foreground">
              You've finished reading this chapter! 
              {readingTime > 0 && ` Reading time: ${formatTime(readingTime)}`}
            </p>
            <div className="flex flex-wrap gap-4 justify-center">
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
                  className="btn-primary"
                >
                  Next Chapter
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              ) : (
                <Button
                  onClick={handleTranslateNext}
                  disabled={isTranslating}
                  className="btn-primary"
                >
                  {isTranslating ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Translating Next...
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
      </div>
      
      {/* Fixed Bottom Navigation */}
      <div className="fixed bottom-0 left-0 right-0 p-4 backdrop-blur-md bg-background/90 border-t">
        <div className="container mx-auto flex justify-between items-center">
          <Button
            variant="outline"
            onClick={() => handleNavigate('prev')}
            disabled={!hasPreviousChapter}
            size="sm"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            <span className="hidden sm:inline">Previous</span>
          </Button>
          
          <div className="flex items-center space-x-4 text-sm text-muted-foreground">
            <span className="hidden sm:inline">{Math.round(readingProgress)}% complete</span>
            <Progress value={readingProgress} className="w-24 h-2" />
          </div>
          
          {hasNextChapter && (
            <Button
              onClick={() => handleNavigate('next')}
              size="sm"
              className="btn-primary"
            >
              <span className="hidden sm:inline">Next</span>
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          )}
        </div>
      </div>
    </div>
  );
};

export default Reader;
