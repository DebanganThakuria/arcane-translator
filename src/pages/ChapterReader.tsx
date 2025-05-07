
import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getNovelById, getChapter } from '../database/db';
import Layout from '../components/Layout';
import Reader from '../components/Reader';
import { useToast } from '@/components/ui/use-toast';
import { Novel, Chapter } from '../types/novel';

const ChapterReader = () => {
  const { novelId, chapterNumber } = useParams<{ novelId: string; chapterNumber: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  const [novel, setNovel] = useState<Novel | undefined>(undefined);
  const [chapter, setChapter] = useState<Chapter | undefined>(undefined);
  const [loading, setLoading] = useState(true);
  
  const chapterNum = chapterNumber ? parseInt(chapterNumber) : 1;
  
  useEffect(() => {
    if (!novelId) {
      navigate('/library');
      return;
    }
    
    try {
      const fetchedNovel = getNovelById(novelId);
      if (fetchedNovel) {
        setNovel(fetchedNovel);
        const fetchedChapter = getChapter(novelId, chapterNum);
        setChapter(fetchedChapter);
        
        if (!fetchedChapter) {
          toast({
            title: "Chapter not found",
            description: "The requested chapter is not available.",
            variant: "destructive",
          });
          navigate(`/novel/${novelId}`);
        }
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
  }, [novelId, chapterNum, navigate, toast]);
  
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
    const targetChapterNum = direction === 'next' ? chapterNum + 1 : chapterNum - 1;
    navigate(`/novel/${novelId}/chapter/${targetChapterNum}`);
    
    // If this was a next chapter and it's not the last one, simulate translation of the next chapter
    if (direction === 'next' && targetChapterNum < novel.chaptersCount - 1) {
      toast({
        title: "Auto-translation",
        description: "Next chapter is being prepared in the background",
      });
    }
  };
  
  const hasPreviousChapter = chapterNum > 1;
  const hasNextChapter = chapterNum < novel.chaptersCount;

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
