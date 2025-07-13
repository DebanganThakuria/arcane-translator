import React, { useEffect, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import Layout from '../components/Layout';
import ChapterList from '../components/ChapterList';
import FirstChapterDialog from '../components/FirstChapterDialog';
import { getNovelById, getChaptersForNovel } from '../database/db';
import { Button } from '@/components/ui/button';
import { ArrowLeft, RefreshCw, Link as LinkIcon, BookOpen, Play } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { formatDistanceToNow } from 'date-fns';
import { Novel, Chapter } from '../types/novel';
import { useToast } from '@/components/ui/use-toast';
import { getReadingProgress, hasReadingProgress, getReadingProgressSummary } from '../utils/readingProgress';
import { refreshNovelWithFallback } from '../services/translationService';

const API_BASE_URL = 'http://localhost:8088';

const NovelDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  const [novel, setNovel] = useState<Novel | undefined>(undefined);
  const [chapters, setChapters] = useState<Chapter[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [firstChapterDialogOpen, setFirstChapterDialogOpen] = useState(false);
  
  const fetchNovelData = async () => {
    if (!id) {
      navigate('/library');
      return;
    }
    
    try {
      setLoading(true);
      const fetchedNovel = await getNovelById(id);
      if (fetchedNovel) {
        setNovel(fetchedNovel);
        const fetchedChapters = await getChaptersForNovel(id);
        setChapters(fetchedChapters);
      } else {
        toast({
          title: "Novel not found",
          description: "The requested novel could not be found.",
          variant: "destructive",
        });
        navigate('/library');
      }
    } catch (error) {
      console.error('Error fetching novel details:', error);
      toast({
        title: "Error",
        description: "Failed to load novel details.",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };
  
  useEffect(() => {
    fetchNovelData();
  }, [id]);
  
  const handleRefresh = async () => {
    if (!id) return;
    
    setRefreshing(true);
    try {
      await refreshNovelWithFallback(id, novel.url);

      toast({
        title: "Refresh Complete",
        variant: "default",
      });

      window.location.reload();
    } catch (error) {
      toast({
        title: "Refresh Failed",
        variant: "destructive",
      });
    } finally {
      setRefreshing(false);
    }
  };
  
  if (loading) {
    return (
      <Layout>
        <div className="container mx-auto py-16 px-4 text-center">
          <p className="text-xl text-muted-foreground mb-4">Loading...</p>
        </div>
      </Layout>
    );
  }
  
  if (!novel) {
    return (
      <Layout>
        <div className="container mx-auto py-16 px-4 text-center">
          <p className="text-xl text-muted-foreground mb-4">Novel not found</p>
          <Button asChild className="gradient-button">
            <Link to="/library">Return to Library</Link>
          </Button>
        </div>
      </Layout>
    );
  }
  
  const continueReading = () => {
    if (!id) return;
    
    // Check localStorage first for reading progress
    const progress = getReadingProgress(id);
    if (progress && progress.lastChapter) {
      const lastReadChapter = chapters.find(chapter => 
        chapter.number === progress.lastChapter
      );
      
      if (lastReadChapter) {
        navigate(`/novel/${id}/chapter/${lastReadChapter.number}`);
        return;
      }
    }
    
    // Fallback to database last read chapter
    if (novel.last_read_chapter_number) {
      const lastReadChapter = chapters.find(chapter => 
        chapter.number === novel.last_read_chapter_number
      );
      
      if (lastReadChapter) {
        navigate(`/novel/${id}/chapter/${lastReadChapter.number}`);
        return;
      }
    }
    
    // Start from first chapter if no progress found
    if (chapters.length > 0) {
      const firstChapter = chapters[0];
      navigate(`/novel/${id}/chapter/${firstChapter.number}`);
    }
  };

  const startReading = () => {
    if (!id || chapters.length === 0) return;
    const firstChapter = chapters[0];
    navigate(`/novel/${id}/chapter/${firstChapter.number}`);
  };

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <div className="flex justify-between items-center mb-6">
          <Button asChild variant="ghost" className="hover:bg-indigo-50/50 hover:text-indigo-600">
            <Link to="/library">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Library
            </Link>
          </Button>
          
          <Button 
            variant="outline" 
            onClick={handleRefresh}
            disabled={refreshing}
            className="border-indigo-200 dark:border-indigo-800 hover:bg-indigo-50 dark:hover:bg-indigo-950/50 text-indigo-600 dark:text-indigo-400"
          >
            <RefreshCw className={`mr-2 h-4 w-4 ${refreshing ? 'animate-spin' : ''}`} />
            {refreshing ? 'Checking...' : 'Check for Updates'}
          </Button>
        </div>
        
        <div className="flex flex-col md:flex-row gap-8">
          <div className="md:w-1/3 lg:w-1/4">
            <div className="glass-card p-6 rounded-lg">
              <div className="book-cover aspect-[2/3] max-w-[300px] mx-auto mb-6 shadow-lg hover:shadow-xl">
                <img 
                  src={novel.cover || '/placeholder.svg'} 
                  alt={novel.title} 
                  className="w-full h-full object-cover rounded-md"
                />
              </div>
              
              <div className="space-y-4">
                {(hasReadingProgress(id || '') || novel.last_read_timestamp) && chapters.length > 0 && (
                  <div className="space-y-2">
                    <Button className="w-full gradient-button" onClick={continueReading}>
                      <BookOpen className="mr-2 h-4 w-4" />
                      Continue Reading
                    </Button>
                    {hasReadingProgress(id || '') && (
                      <p className="text-xs text-center text-muted-foreground">
                        {getReadingProgressSummary(id || '')}
                      </p>
                    )}
                  </div>
                )}
                
                {!hasReadingProgress(id || '') && !novel.last_read_timestamp && chapters.length > 0 && (
                  <Button className="w-full gradient-button" onClick={startReading}>
                    <Play className="mr-2 h-4 w-4" />
                    Start Reading
                  </Button>
                )}
                
                {hasReadingProgress(id || '') && chapters.length > 0 && (
                  <Button 
                    variant="outline" 
                    className="w-full border-indigo-200 hover:bg-indigo-50 text-indigo-600"
                    onClick={startReading}
                  >
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Start from Beginning
                  </Button>
                )}
                
                <div className="glass-card p-4 rounded-md text-sm">
                  <p className="mb-1"><span className="text-indigo-600 font-medium">Status:</span> {novel.status}</p>
                  <p className="mb-1"><span className="text-indigo-600 font-medium">Chapters:</span> {novel.chapters_count}</p>
                  <p className="mb-1">
                    <span className="text-indigo-600 font-medium">Last Updated:</span>{" "}
                    {novel.last_updated ? formatDistanceToNow(new Date(parseInt(String(novel.last_updated)) * 1000)) : "unknown"} ago
                  </p>
                  {novel.source && (
                    <p className="mb-1">
                      <span className="text-indigo-600 font-medium">Source:</span>{' '}
                      <span>{novel.source}</span>
                    </p>
                  )}
                </div>
                
                {novel.genres && novel.genres.length > 0 && (
                  <div className="glass-card p-4 rounded-md">
                    <p className="text-sm text-indigo-600 font-medium mb-2">Genres:</p>
                    <div className="flex flex-wrap gap-2">
                      {novel.genres.map(genre => (
                        <Link to={`/genre/${genre}`} key={genre}>
                          <Badge 
                            variant="secondary" 
                            className="bg-gradient-to-r from-indigo-600/20 to-blue-400/20 dark:from-indigo-500/30 dark:to-blue-500/30 text-indigo-700 dark:text-indigo-200 hover:from-indigo-600/30 hover:to-blue-400/30 dark:hover:from-indigo-500/40 dark:hover:to-blue-500/40 cursor-pointer transition-colors"
                          >
                            {genre}
                          </Badge>
                        </Link>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
          
          <div className="md:w-2/3 lg:w-3/4">
            <div className="glass-card p-6 rounded-lg mb-6">
              <h1 className="text-3xl font-bold mb-2 text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">{novel.title}</h1>
              
              {novel.original_title && (
                <span>
                  <p className="text-lg text-muted-foreground mb-2">{novel.original_title} | <a href={novel.url} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:text-blue-600">{novel.url}</a></p>
                </span>
              )}
              
              {novel.author && (
                <p className="mb-4">by <span className="font-medium text-indigo-600">{novel.author}</span></p>
              )}
              
              <div 
                className="mb-0 prose prose-sm max-w-none dark:prose-invert"
                style={{
                  color: 'inherit',
                  ['--tw-prose-body' as string]: 'var(--foreground)',
                  ['--tw-prose-headings' as string]: 'var(--foreground)',
                  ['--tw-prose-links' as string]: 'var(--primary)',
                  ['--tw-prose-bold' as string]: 'var(--foreground)',
                  ['--tw-prose-counters' as string]: 'var(--muted-foreground)',
                  ['--tw-prose-bullets' as string]: 'var(--muted-foreground)',
                } as React.CSSProperties}
                dangerouslySetInnerHTML={{ __html: novel.summary }}
              />
            </div>
            
            <div className="glass-card p-6 rounded-lg">
              <h2 className="text-2xl font-bold mb-4 text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">Chapters</h2>
              <Separator className="my-4 bg-gradient-to-r from-indigo-100 to-blue-100" />
              <ChapterList chapters={chapters} novelId={id} />
            </div>
          </div>
        </div>
        
        {chapters.length === 0 && (
          <div className="fixed bottom-4 right-4">
            <Button 
              onClick={() => setFirstChapterDialogOpen(true)} 
              className="gradient-button flex items-center"
            >
              <LinkIcon className="mr-2 h-5 w-5" />
              First Chapter
            </Button>
          </div>
        )}
        
        {chapters.length === 0 && (
          <FirstChapterDialog 
            isOpen={firstChapterDialogOpen} 
            onOpenChange={setFirstChapterDialogOpen}
            novelId={id || ''}
            novelUrl={novel.url}
            onSuccess={() => {
              toast({
                title: "Ready to Read",
                description: "You can now start reading the first chapter.",
              });
              navigate(`/novel/${id}/chapter/1`);
            }}
          />
        )}
      </div>
    </Layout>
  );
};

export default NovelDetail;
