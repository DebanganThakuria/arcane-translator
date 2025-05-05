
import React, { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getNovelById, getChapter } from '../data/mockData';
import Layout from '../components/Layout';
import Reader from '../components/Reader';
import { useToast } from '@/components/ui/use-toast';

const ChapterReader = () => {
  const { novelId, chapterNumber } = useParams<{ novelId: string; chapterNumber: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  
  const chapterNum = chapterNumber ? parseInt(chapterNumber) : 1;
  
  if (!novelId) {
    navigate('/library');
    return null;
  }
  
  const novel = getNovelById(novelId);
  const chapter = getChapter(novelId, chapterNum);
  
  useEffect(() => {
    if (!chapter && novel) {
      toast({
        title: "Chapter not found",
        description: "The requested chapter is not available.",
        variant: "destructive",
      });
      navigate(`/novel/${novelId}`);
    }
  }, [chapter, novel, navigate, novelId, toast]);
  
  if (!novel || !chapter) {
    return null;
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
