import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getNovelById, getChapterByNumber, getChaptersForNovel } from '../database/db';
import Layout from '../components/Layout';
import Reader from '../components/Reader';
import { useToast } from '@/components/ui/use-toast';
import { Novel, Chapter } from '../types/novel';
import { saveReadingProgress } from '../utils/readingProgress';

const ChapterReader = () => {
  const { novelId, chapterNumber } = useParams<{ novelId: string; chapterNumber: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  const [novel, setNovel] = useState<Novel | undefined>(undefined);
  const [chapter, setChapter] = useState<Chapter | undefined>(undefined);
  const [chapters, setChapters] = useState<Chapter[]>([]);
  const [loading, setLoading] = useState(true);
  const [hasNextChapter, setHasNextChapter] = useState(false);

  const chapterNum = chapterNumber ? parseInt(chapterNumber) : 1;
  const hasPreviousChapter = chapterNum > 1;
  
  useEffect(() => {
    const fetchData = async () => {
      if (!novelId) {
        navigate('/library');
        return;
      }
      
      try {
        setLoading(true);
        
        const fetchedNovel = await getNovelById(novelId);
        
        if (fetchedNovel) {
          setNovel(fetchedNovel);
          
          // Fetch all chapters for navigation
          const allChapters = await getChaptersForNovel(novelId);
          setChapters(allChapters);
          
          const fetchedChapter = await getChapterByNumber(novelId, chapterNum);
          
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
          
          // Save reading progress when chapter loads
          saveReadingProgress(novelId, chapterNum, 0, fetchedChapter.title);
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

    fetchData();
  }, [novelId, chapterNum, navigate, toast]);

  useEffect(() => {
    const checkNextChapter = async () => {
      if (!novelId) return;
      const nextChapter = await getChapterByNumber(novelId, chapterNum + 1);
      setHasNextChapter(!!nextChapter);
    };

    checkNextChapter();
  }, [chapterNum, novelId]);

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
    
    if (targetChapterNum < 1 || (novel.chapters_count && targetChapterNum > novel.chapters_count)) {
      return;
    }
    
    navigate(`/novel/${novelId}/chapter/${targetChapterNum}`);
  };

  const handleTranslateNext = async () => {
    if (!novelId) return;
    
    try {
      const nextChapter = await getChapterByNumber(novelId, chapterNum + 1);
      if (nextChapter) {
        setHasNextChapter(true);
      }
    } catch (error) {
      console.error('Error checking for next chapter after translation:', error);
    }
  };

  return (
    <Layout hideNavigation>
      <Reader
        chapter={chapter}
        novel={novel}
        hasPreviousChapter={hasPreviousChapter}
        hasNextChapter={hasNextChapter}
        onNavigate={handleNavigate}
        onTranslateNext={handleTranslateNext}
      />
    </Layout>
  );
};

export default ChapterReader;
