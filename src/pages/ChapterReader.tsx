import React, { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getNovelById, getChapter, saveChapter } from '../database/db';
import Layout from '../components/Layout';
import Reader from '../components/Reader';
import { useToast } from '@/components/ui/use-toast';
import { Novel, Chapter } from '../types/novel';
import { getNovel, translateChapter, getChapterUrl } from '../services/translationService';

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
          
          // Try to get the chapter from local database first
          let fetchedChapter = await getChapter(novelId, chapterNum);
          
          // If not found locally, translate it from the backend
          if (!fetchedChapter) {
            try {
              toast({
                title: "Translating",
                description: `Chapter ${chapterNum} is being translated...`,
              });
              
              let chapterUrl;
              // If we have a firstChapterUrl and urlPattern, use those to generate the URL
              if (fetchedNovel.firstChapterUrl && fetchedNovel.urlPattern) {
                try {
                  chapterUrl = await getChapterUrl(novelId, chapterNum);
                } catch (urlError) {
                  console.error('Error generating chapter URL:', urlError);
                  // Fall back to default URL format
                  chapterUrl = `${fetchedNovel.url}/chapter-${chapterNum}`;
                }
              } else {
                // Use default URL format
                chapterUrl = `${fetchedNovel.url}/chapter-${chapterNum}`;
              }
              
              // Try to translate and fetch the chapter from the source
              const translatedChapter = await translateChapter(
                novelId,
                chapterNum,
                chapterUrl,
                undefined,
                true // Save to database
              );
              
              // Create a chapter object from the translated content
              if (translatedChapter) {
                fetchedChapter = {
                  id: `${novelId}-${chapterNum}`, // Generate a temporary ID
                  novelId: novelId,
                  number: chapterNum,
                  title: translatedChapter.translated_chapter_title,
                  content: translatedChapter.translated_chapter_contents,
                  dateTranslated: Date.now(),
                  wordCount: translatedChapter.word_count
                };
                
                // Save the chapter to the database
                await saveChapter(fetchedChapter);
              }
            } catch (translationError) {
              console.error('Error translating chapter:', translationError);
              // Continue with local fetch attempt
            }
          }
          
          setChapter(fetchedChapter);
          
          if (!fetchedChapter) {
            toast({
              title: "Chapter not found",
              description: "The requested chapter is not available.",
              variant: "destructive",
            });
            navigate(`/novel/${novelId}`);
            return;
          }
          
          // Start background translation for next chapter if we're not at the last chapter
          const nextChapterNum = chapterNum + 1;
          if (nextChapterNum <= fetchedNovel.chaptersCount) {
            backgroundTranslateNextChapter(novelId, nextChapterNum, fetchedNovel);
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
    };

    fetchData();
  }, [novelId, chapterNum, navigate, toast]);
  
  // Function to translate the next chapter in the background
  const backgroundTranslateNextChapter = async (novelId: string, chapterNumber: number, novel: Novel) => {
    // Check if background translation is already in progress for this chapter
    if (
      backgroundTranslationRef.current.inProgress && 
      backgroundTranslationRef.current.chapterNumber === chapterNumber
    ) {
      return; // Already translating this chapter
    }
    
    // Check if chapter already exists in database
    try {
      const existingChapter = await getChapter(novelId, chapterNumber);
      if (existingChapter) {
        // Chapter already exists, no need to translate
        return;
      }
      
      // Set background translation status
      backgroundTranslationRef.current = {
        inProgress: true,
        chapterNumber: chapterNumber
      };
      
      // Try to get the URL for the next chapter in the following ways:
      // 1. Check if current chapter has a nextChapterUrl property
      // 2. Generate URL using the pattern if available
      // 3. Fall back to default URL format
      let chapterUrl;
      
      // Check if the current chapter has nextChapterUrl
      const currentChapterNum = chapterNumber - 1;
      try {
        const currentChapter = await getChapter(novelId, currentChapterNum);
        if (currentChapter && currentChapter.nextChapterUrl) {
          chapterUrl = currentChapter.nextChapterUrl;
          console.log(`Using nextChapterUrl from chapter ${currentChapterNum}: ${chapterUrl}`);
        }
      } catch (error) {
        console.error('Error fetching previous chapter for nextChapterUrl:', error);
      }
      
      // If nextChapterUrl is not available, generate one using the pattern
      if (!chapterUrl && novel.firstChapterUrl && novel.urlPattern) {
        try {
          chapterUrl = await getChapterUrl(novelId, chapterNumber);
        } catch (urlError) {
          console.error('Error generating background chapter URL:', urlError);
          // Fall back to default URL format
          chapterUrl = `${novel.url}/chapter-${chapterNumber}`;
        }
      } else if (!chapterUrl) {
        // Use default URL format as last resort
        chapterUrl = `${novel.url}/chapter-${chapterNumber}`;
      }
      
      // Translate the chapter silently in the background
      const translatedChapter = await translateChapter(
        novelId,
        chapterNumber,
        chapterUrl,
        undefined,
        true // Save to database
      );
      
      // Create a chapter object from the translated content
      if (translatedChapter) {
        const chapter = {
          id: `${novelId}-${chapterNumber}`,
          novelId: novelId,
          number: chapterNumber,
          title: translatedChapter.translated_chapter_title,
          content: translatedChapter.translated_chapter_contents,
          dateTranslated: Date.now(),
          wordCount: translatedChapter.word_count
        };
        
        // Save the chapter to database
        await saveChapter(chapter);
        console.log(`Chapter ${chapterNumber} background translation complete`);
      }
    } catch (error) {
      console.error(`Error in background translation for chapter ${chapterNumber}:`, error);
    } finally {
      // Reset background translation status
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
    const targetChapterNum = direction === 'next' ? chapterNum + 1 : chapterNum - 1;
    
    // If we're going to the next chapter and have nextChapterUrl, we can use it for analytics
    // or potentially preload content in the future
    if (direction === 'next' && chapter.nextChapterUrl) {
      console.log(`Next chapter URL available: ${chapter.nextChapterUrl}`);
    }
    
    navigate(`/novel/${novelId}/chapter/${targetChapterNum}`);
    
    // If this was a next chapter and it's not the last one, start background translation of the next-next chapter
    if (direction === 'next' && targetChapterNum < novel.chaptersCount - 1 && novel) {
      // Start background translation for the next-next chapter (looking ahead)
      const nextNextChapterNum = targetChapterNum + 1;
      
      // Show toast with slightly longer duration
      toast({
        title: "Auto-translation",
        description: "Next chapter is being prepared in the background",
        duration: 3000,
      });
      
      // Use a slight delay to avoid too many simultaneous requests
      setTimeout(() => {
        backgroundTranslateNextChapter(novelId, nextNextChapterNum, novel);
      }, 1000);
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
