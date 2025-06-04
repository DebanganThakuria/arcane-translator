import React, { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {getNovelById, getChapterByNumber} from '../database/db';
import Layout from '../components/Layout';
import Reader from '../components/Reader';
import { useToast } from '@/components/ui/use-toast';
import { Progress } from '@/components/ui/progress';
import { Loader2 } from 'lucide-react';
import { Novel, Chapter } from '../types/novel';
import { getNovel } from '../services/translationService';

const API_BASE_URL = 'http://localhost:8088';

const ChapterReader = () => {
  const { novelId, chapterNumber } = useParams<{ novelId: string; chapterNumber: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  const [novel, setNovel] = useState<Novel | undefined>(undefined);
  const [chapter, setChapter] = useState<Chapter | undefined>(undefined);
  const [loading, setLoading] = useState(true);
  const [translationProgress, setTranslationProgress] = useState<{
    inProgress: boolean;
    chapterNumber: number | null;
    progress: number;
  }>({ inProgress: false, chapterNumber: null, progress: 0 });

  const [translationState, setTranslationState] = useState<{
    inProgress: boolean;
    chapterNumber: number | null;
    novelId: string | null;
    progress: number;
  }>({ inProgress: false, chapterNumber: null, novelId: null, progress: 0 });
  
  const chapterNum = chapterNumber ? parseInt(chapterNumber) : 1;
  
  useEffect(() => {
    const fetchData = async () => {
      if (!novelId) {
        navigate('/library');
        return;
      }
      
      try {
        setLoading(true);
        
        // Try to get the novel from local database first
        let fetchedNovel = await getNovelById(novelId);
        
        // If not found locally, try to fetch from backend API
        if (!fetchedNovel) {
          try {
            fetchedNovel = await getNovel(novelId);
          } catch (apiError) {
            console.error('API error fetching novel:', apiError);
            // Continue with local fetch attempt
          }
        }
        
        if (fetchedNovel) {
          setNovel(fetchedNovel);
          
          // Try to get the chapter from local database
          const fetchedChapter = await getChapterByNumber(novelId, chapterNum);
          
          // If chapter doesn't exist, redirect to novel page
          if (!fetchedChapter) {
            toast({
              title: "Chapter not available",
              description: `Chapter ${chapterNum} is not available.`,
              variant: "destructive",
            });
            navigate(`/novel/${novelId}`);
            return;
          }

          setChapter(fetchedChapter);

          // Ensure next chapter is available in background
          ensureNextChapterAvailable(chapterNum);
        } else {
          toast({
            title: "Novel not found",
            description: "The requested novel is not available.",
            variant: "destructive",
          });
          navigate('/library');
        }
      } catch (error) {
        console.error('Error fetching chapter:', error);
        toast({
          title: "Error",
          description: "Failed to load chapter.",
          variant: "destructive",
        });
        navigate('/library');
      } finally {
        setLoading(false);
      }
    };

    fetchData()
    
    // Cleanup function to reset translation state
    return () => {
      setTranslationState({
        inProgress: false,
        chapterNumber: null,
        novelId: null,
        progress: 0
      });
    };
  }, [novelId, chapterNum, navigate, toast]);
  
  // Effect to ensure next chapter is prepared when current chapter is loaded
  useEffect(() => {
    if (chapter && novel) {
      ensureNextChapterAvailable(chapterNum);
    }
  }, [chapter, chapterNum, novel]);
  
  // Function to ensure next chapter is available
  const ensureNextChapterAvailable = async (currentChapterNum: number) => {
    if (!novelId || !novel) return;
    
    const nextChapterNum = currentChapterNum + 1;
    if (nextChapterNum > (novel.chapters_count || 0)) return; // No more chapters
    
    // Skip if already translating this chapter
    if (translationState.inProgress) {
      return;
    }

    try {
      // Check if next chapter already exists
      const existingChapter = await getChapterByNumber(novelId, nextChapterNum);
      if (existingChapter) return; // Chapter already exists
      
      console.log(`Requesting translation for chapter ${nextChapterNum} in background...`);
      
      // Set translation in progress
      const newState = {
        inProgress: true,
        chapterNumber: nextChapterNum,
        novelId: novelId,
        progress: 0
      };
      setTranslationState(newState);
      setTranslationProgress({
        inProgress: true,
        chapterNumber: nextChapterNum,
        progress: 0
      });
      
      // Start progress simulation
      const progressInterval = setInterval(() => {
        setTranslationProgress(prev => ({
          ...prev,
          progress: Math.min(prev.progress + 5, 90) // Cap at 90% until complete
        }));
      }, 1000);
      
      // Don't await this, let it happen in background
      fetch(`${API_BASE_URL}/novels/translate/chapter`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          novel_id:novelId,
        })
      })
      .then(() => {
        clearInterval(progressInterval);
        setTranslationProgress(prev => ({ ...prev, progress: 100 }));
        toast({
          title: "Translation Complete",
          description: `Chapter ${nextChapterNum} has been translated and is ready to read!`,
          variant: "default",
        });
        // Reset after showing completion
        setTimeout(() => {
          setTranslationProgress({ inProgress: false, chapterNumber: null, progress: 0 });
        }, 2000);
      })
      .catch(error => {
        clearInterval(progressInterval);
        console.error(`Background translation failed for chapter ${nextChapterNum}:`, error);
        toast({
          title: "Translation Failed",
          description: `Failed to translate chapter ${nextChapterNum}. Please try again later.`,
          variant: "destructive",
        });
      });
      
    } catch (error) {
      console.error('Error translating next chapter:', error);
    } finally {
      // Reset translation status
      setTranslationState({
        inProgress: false,
        chapterNumber: null,
        novelId: null,
        progress: 0
      });
    }
  };
  
  if (loading || !novel || !chapter) {
    return (
      <Layout hideNavigation>
        <div className="container mx-auto py-16 px-4 text-center">
          <p className="text-xl text-muted-foreground mb-4">Loading...</p>
        </div>
      </Layout>
    );
  }

  const handleNavigate = (direction: 'next' | 'prev') => {
    if (!novel) return;
    
    const targetChapterNum = direction === 'next' ? chapterNum + 1 : chapterNum - 1;
    
    // Validate chapter number
    if (targetChapterNum < 1 || (novel.chapters_count && targetChapterNum > novel.chapters_count)) {
      return;
    }
    
    // Navigate to the target chapter
    navigate(`/novel/${novelId}/chapter/${targetChapterNum}`);
  };
  
  const hasPreviousChapter = chapterNum > 1;
  const hasNextChapter = chapterNum < novel.chapters_count;

  return (
    <Layout hideNavigation>
      <div className="relative
        {translationProgress.inProgress ? 'pb-16' : ''}
      ">
        <Reader
          chapter={chapter}
          novel={novel}
          hasPreviousChapter={hasPreviousChapter}
          hasNextChapter={hasNextChapter}
          onNavigate={handleNavigate}
        />
        
        {translationProgress.inProgress && (
          <div className="fixed bottom-0 left-0 right-0 bg-background/90 backdrop-blur-sm p-4 border-t border-border shadow-lg">
            <div className="container mx-auto max-w-4xl">
              <div className="flex items-center gap-4">
                <Loader2 className="h-5 w-5 animate-spin text-primary" />
                <div className="flex-1">
                  <div className="text-sm font-medium mb-1">
                    Translating Chapter {translationProgress.chapterNumber}...
                  </div>
                  <Progress 
                    value={translationProgress.progress} 
                    className="h-2"
                  />
                  <div className="text-xs text-muted-foreground mt-1 text-right">
                    {translationProgress.progress}%
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
};

export default ChapterReader;
