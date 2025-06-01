import React, { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {getNovelById, getChapterByNumber} from '../database/db';
import Layout from '../components/Layout';
import Reader from '../components/Reader';
import { useToast } from '@/components/ui/use-toast';
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
  
  // Ref to track background translation status to avoid duplicate translations
  const backgroundTranslationRef = useRef<{
    inProgress: boolean;
    chapterNumber: number | null;
  }>({ inProgress: false, chapterNumber: null });
  
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
  }, [novelId, chapterNum, navigate, toast]);
  
  // Function to ensure next chapter is available
  const ensureNextChapterAvailable = async (currentChapterNum: number) => {
    if (!novelId || !novel) return;
    
    const nextChapterNum = currentChapterNum + 1;
    if (nextChapterNum > (novel.chapters_count || 0)) return; // No more chapters
    
    // Skip if already translating this chapter
    if (backgroundTranslationRef.current.inProgress && backgroundTranslationRef.current.chapterNumber === nextChapterNum) {
      return;
    }
    
    try {
      // Check if next chapter already exists
      const existingChapter = await getChapterByNumber(novelId, nextChapterNum);
      if (existingChapter) return; // Chapter already exists
      
      console.log(`Requesting translation for chapter ${nextChapterNum} in background...`);
      
      // Set translation in progress
      backgroundTranslationRef.current = {
        inProgress: true,
        chapterNumber: nextChapterNum
      };
      
      // Don't await this, let it happen in background
      fetch(`${API_BASE_URL}/novels/translate/chapter`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          novelId,
        })
      }).catch(error => {
        console.error(`Background translation failed for chapter ${nextChapterNum}:`, error);
      });
      
    } catch (error) {
      console.error('Error translating next chapter:', error);
    } finally {
      // Reset translation status
      backgroundTranslationRef.current = {
        inProgress: false,
        chapterNumber: null
      };
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
      <Reader
        chapter={chapter}
        novel={novel}
        hasPreviousChapter={hasPreviousChapter}
        hasNextChapter={hasNextChapter}
        onNavigate={handleNavigate}
      />
    </Layout>
  );
};

export default ChapterReader;
