import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Chapter, Novel } from '../types/novel';
import { ArrowLeft, ArrowRight, Loader2, Settings, BookOpen, Clock, Eye, ScrollText, Maximize, Minimize } from 'lucide-react';
import { translateChapter } from '../services/translationService';
import { useToast } from '@/components/ui/use-toast';
import { Progress } from '@/components/ui/progress';
import { Slider } from '@/components/ui/slider';
import { saveReadingProgress } from '../utils/readingProgress';
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
  const [showSettings, setShowSettings] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [showFullscreenBanner, setShowFullscreenBanner] = useState(false);
  const navigate = useNavigate();
  const { toast } = useToast();
  const contentRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Scroll to top when chapter changes
    window.scrollTo(0, 0);
    setReadingProgress(0);

    // Load saved preferences
    const savedTheme = localStorage.getItem('readerTheme') as 'light' | 'sepia' | 'dark';
    const savedFontSize = localStorage.getItem('readerFontSize');
    const savedLineHeight = localStorage.getItem('readerLineHeight');
    
    if (savedTheme) setTheme(savedTheme);
    if (savedFontSize) setFontSize(parseInt(savedFontSize));
    if (savedLineHeight) setLineHeight(parseFloat(savedLineHeight));
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

    window.addEventListener('scroll', handleScroll);
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, []);

  // Save reading progress to localStorage
  useEffect(() => {
    const saveProgressDebounced = setTimeout(() => {
      if (novel.id && chapter.number) {
        saveReadingProgress(
          novel.id,
          chapter.number,
          readingProgress,
          chapter.title
        );
      }
    }, 1000); // Debounce for 1 second

    return () => clearTimeout(saveProgressDebounced);
  }, [readingProgress, novel.id, chapter.number, chapter.title]);

  // Save progress when user leaves the page
  useEffect(() => {
    const handleBeforeUnload = () => {
      if (novel.id && chapter.number) {
        saveReadingProgress(
          novel.id,
          chapter.number,
          readingProgress,
          chapter.title
        );
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
    };
  }, [novel.id, chapter.number, readingProgress, chapter.title]);

  // Fullscreen and ESC key handler
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape' && isFullscreen) {
        exitFullscreen();
      }
    };

    const handleFullscreenChange = () => {
      const isNowFullscreen = !!document.fullscreenElement;
      setIsFullscreen(isNowFullscreen);
      if (!isNowFullscreen) {
        setShowFullscreenBanner(false);
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    document.addEventListener('fullscreenchange', handleFullscreenChange);

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      document.removeEventListener('fullscreenchange', handleFullscreenChange);
    };
  }, [isFullscreen]);

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

  const enterFullscreen = async () => {
    try {
      await document.documentElement.requestFullscreen();
      setIsFullscreen(true);
      setShowFullscreenBanner(true);
      // Hide banner after 2 seconds
      setTimeout(() => setShowFullscreenBanner(false), 2000);
    } catch (error) {
      console.error('Error entering fullscreen:', error);
      toast({
        title: "Error",
        description: "Failed to enter fullscreen mode",
        variant: "destructive",
      });
    }
  };

  const exitFullscreen = async () => {
    try {
      if (document.fullscreenElement) {
        await document.exitFullscreen();
      }
      setIsFullscreen(false);
      setShowFullscreenBanner(false);
    } catch (error) {
      console.error('Error exiting fullscreen:', error);
    }
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
      {/* Fullscreen Banner */}
      {showFullscreenBanner && (
        <div className="fixed inset-0 z-50 flex items-center justify-center pointer-events-none">
          <div className="bg-black/80 text-white px-6 py-3 rounded-lg backdrop-blur-sm">
            <p className="text-center text-sm font-medium">
              Press <kbd className="px-2 py-1 bg-white/20 rounded text-xs">ESC</kbd> to toggle fullscreen
            </p>
          </div>
        </div>
      )}

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
                  <ScrollText className="h-4 w-4" />
                  <span>{novel.title}</span>
                </div>
                <div className="flex items-center space-x-1">
                  <BookOpen className="h-4 w-4" />
                  <span>Chapter {chapter.number}</span>
                </div>
                <div className="flex items-center space-x-1">
                  <Clock className="h-4 w-4" />
                  <span>~{estimatedReadingTime} min read</span>
                </div>
                <div className="flex items-center space-x-1">
                  <Eye className="h-4 w-4" />
                  <span>{Math.round(readingProgress)}%</span>
                </div>
              </div>
            </div>

            <div className="flex items-center gap-2">
              {/* Fullscreen Button */}
              <Button
                variant="outline"
                size="sm"
                onClick={isFullscreen ? exitFullscreen : enterFullscreen}
                className={`gap-2 ${
                  theme === 'dark' 
                    ? 'bg-gray-800 border-gray-600 text-gray-100 hover:bg-gray-700 hover:text-white' 
                    : theme === 'sepia'
                    ? 'bg-amber-50 border-amber-300 text-amber-800 hover:bg-amber-100'
                    : 'bg-white border-gray-300 text-gray-700 hover:bg-gray-50'
                }`}
              >
                {isFullscreen ? (
                  <>
                    <Minimize className="h-4 w-4" />
                    <span className="hidden sm:inline">Exit Fullscreen</span>
                  </>
                ) : (
                  <>
                    <Maximize className="h-4 w-4" />
                    <span className="hidden sm:inline">Fullscreen</span>
                  </>
                )}
              </Button>

              {/* Reader Settings */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button 
                    variant="outline" 
                    size="sm" 
                    className={`gap-2 ${
                      theme === 'dark' 
                        ? 'bg-gray-800 border-gray-600 text-gray-100 hover:bg-gray-700 hover:text-white' 
                        : theme === 'sepia'
                        ? 'bg-amber-50 border-amber-300 text-amber-800 hover:bg-amber-100'
                        : 'bg-white border-gray-300 text-gray-700 hover:bg-gray-50'
                    }`}
                  >
                    <Settings className="h-4 w-4" />
                    <span className="hidden sm:inline">Settings</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-96 p-0">
                  <div className="p-6">
                    <h3 className="text-lg font-semibold mb-6 text-center">Reading Settings</h3>
                    
                    {/* Theme Selection */}
                    <div className="mb-8">
                      <label className="text-sm font-medium mb-4 block text-muted-foreground">Reading Theme</label>
                      <div className="grid grid-cols-3 gap-3">
                        <button
                          onClick={() => handleThemeChange('light')}
                          className={`relative p-4 rounded-xl border-2 transition-all hover:scale-105 ${
                            theme === 'light' 
                              ? 'border-primary ring-2 ring-primary/20 shadow-lg' 
                              : 'border-gray-200 hover:border-gray-300'
                          }`}
                        >
                          <div className="w-full h-16 bg-white rounded-lg mb-3 shadow-inner border border-gray-100 flex items-center justify-center">
                            <div className="text-xs font-medium text-gray-800">Aa</div>
                          </div>
                          <div className="text-xs font-medium text-center">Light</div>
                          {theme === 'light' && (
                            <div className="absolute -top-1 -right-1 w-6 h-6 bg-primary rounded-full flex items-center justify-center">
                              <div className="w-2 h-2 bg-white rounded-full"></div>
                            </div>
                          )}
                        </button>
                        
                        <button
                          onClick={() => handleThemeChange('sepia')}
                          className={`relative p-4 rounded-xl border-2 transition-all hover:scale-105 ${
                            theme === 'sepia' 
                              ? 'border-primary ring-2 ring-primary/20 shadow-lg' 
                              : 'border-gray-200 hover:border-gray-300'
                          }`}
                        >
                          <div className="w-full h-16 bg-amber-50 rounded-lg mb-3 shadow-inner border border-amber-100 flex items-center justify-center">
                            <div className="text-xs font-medium text-amber-800">Aa</div>
                          </div>
                          <div className="text-xs font-medium text-center">Sepia</div>
                          {theme === 'sepia' && (
                            <div className="absolute -top-1 -right-1 w-6 h-6 bg-primary rounded-full flex items-center justify-center">
                              <div className="w-2 h-2 bg-white rounded-full"></div>
                            </div>
                          )}
                        </button>
                        
                        <button
                          onClick={() => handleThemeChange('dark')}
                          className={`relative p-4 rounded-xl border-2 transition-all hover:scale-105 ${
                            theme === 'dark' 
                              ? 'border-primary ring-2 ring-primary/20 shadow-lg' 
                              : 'border-gray-200 hover:border-gray-300'
                          }`}
                        >
                          <div className="w-full h-16 bg-gray-800 rounded-lg mb-3 shadow-inner border border-gray-700 flex items-center justify-center">
                            <div className="text-xs font-medium text-gray-100">Aa</div>
                          </div>
                          <div className="text-xs font-medium text-center">Dark</div>
                          {theme === 'dark' && (
                            <div className="absolute -top-1 -right-1 w-6 h-6 bg-primary rounded-full flex items-center justify-center">
                              <div className="w-2 h-2 bg-white rounded-full"></div>
                            </div>
                          )}
                        </button>
                      </div>
                    </div>

                    {/* Typography Controls */}
                    <div className="space-y-6">
                      {/* Font Size */}
                      <div>
                        <div className="flex items-center justify-between mb-3">
                          <label className="text-sm font-medium text-muted-foreground">Font Size</label>
                          <div className="bg-primary/10 text-primary px-2 py-1 rounded text-sm font-medium">
                            {fontSize}px
                          </div>
                        </div>
                        <div className="px-3">
                          <Slider
                            value={[fontSize]}
                            onValueChange={handleFontSizeChange}
                            max={28}
                            min={12}
                            step={1}
                            className="w-full"
                          />
                          <div className="flex justify-between text-xs text-muted-foreground mt-2">
                            <span>Small</span>
                            <span>Large</span>
                          </div>
                        </div>
                      </div>

                      {/* Line Height */}
                      <div>
                        <div className="flex items-center justify-between mb-3">
                          <label className="text-sm font-medium text-muted-foreground">Line Height</label>
                          <div className="bg-primary/10 text-primary px-2 py-1 rounded text-sm font-medium">
                            {lineHeight.toFixed(1)}
                          </div>
                        </div>
                        <div className="px-3">
                          <Slider
                            value={[lineHeight]}
                            onValueChange={handleLineHeightChange}
                            max={2.5}
                            min={1.2}
                            step={0.1}
                            className="w-full"
                          />
                          <div className="flex justify-between text-xs text-muted-foreground mt-2">
                            <span>Tight</span>
                            <span>Loose</span>
                          </div>
                        </div>
                      </div>
                    </div>
                    
                    {/* Preview */}
                    <div className="mt-6 pt-6 border-t border-gray-200">
                      <div className="text-xs text-muted-foreground mb-2">Preview</div>
                      <div 
                        className={`p-4 rounded-lg border ${
                          theme === 'light' ? 'bg-white border-gray-200' :
                          theme === 'sepia' ? 'bg-amber-50 border-amber-200' :
                          'bg-gray-800 border-gray-700 text-gray-100'
                        }`}
                        style={{ 
                          fontSize: `${fontSize}px`,
                          lineHeight: lineHeight
                        }}
                      >
                        The quick brown fox jumps over the lazy dog. This is how your reading experience will look.
                      </div>
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
              <span>~{estimatedReadingTime} min read</span>
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
    </div>
  );
};

export default Reader;
