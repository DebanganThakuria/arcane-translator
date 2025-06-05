
import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getNovelById, getChapterByNumber } from '../database/db';
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
  
  const hasPreviousChapter = chapterNum > 1;
  const hasNextChapter = chapterNum < (novel.chapters_count || 0);

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
