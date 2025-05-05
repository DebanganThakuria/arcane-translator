
import React from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import Layout from '../components/Layout';
import ChapterList from '../components/ChapterList';
import { getNovelById, getChaptersForNovel } from '../data/mockData';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { formatDistanceToNow } from 'date-fns';

const NovelDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  
  if (!id) {
    navigate('/library');
    return null;
  }
  
  const novel = getNovelById(id);
  const chapters = getChaptersForNovel(id);
  
  if (!novel) {
    return (
      <Layout>
        <div className="container mx-auto py-16 px-4 text-center">
          <p className="text-xl text-muted-foreground mb-4">Novel not found</p>
          <Button asChild>
            <Link to="/library">Return to Library</Link>
          </Button>
        </div>
      </Layout>
    );
  }
  
  const continueReading = () => {
    if (novel.lastRead) {
      navigate(`/novel/${id}/chapter/${novel.lastRead.chapterId.split('-')[1]}`);
    } else if (chapters.length > 0) {
      navigate(`/novel/${id}/chapter/1`);
    }
  };

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <Button asChild variant="ghost" className="mb-6">
          <Link to="/library">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Library
          </Link>
        </Button>
        
        <div className="flex flex-col md:flex-row gap-8">
          {/* Left Column - Novel Info */}
          <div className="md:w-1/3 lg:w-1/4">
            <div className="book-cover aspect-[2/3] max-w-[300px] mx-auto mb-6">
              <img 
                src={novel.cover || '/placeholder.svg'} 
                alt={novel.title} 
                className="w-full h-full object-cover rounded-md shadow-lg"
              />
            </div>
            
            <div className="space-y-4">
              {novel.lastRead && (
                <Button className="w-full" onClick={continueReading}>
                  Continue Reading
                </Button>
              )}
              
              {!novel.lastRead && chapters.length > 0 && (
                <Button className="w-full" onClick={continueReading}>
                  Start Reading
                </Button>
              )}
              
              <div className="text-sm">
                <p className="mb-1"><span className="text-muted-foreground">Status:</span> {novel.status}</p>
                <p className="mb-1"><span className="text-muted-foreground">Chapters:</span> {novel.chaptersCount}</p>
                <p className="mb-1">
                  <span className="text-muted-foreground">Last Updated:</span> {formatDistanceToNow(novel.lastUpdated)} ago
                </p>
                {novel.source && (
                  <p className="mb-1">
                    <span className="text-muted-foreground">Source:</span>{' '}
                    <Link to={`/source/${novel.source}`} className="hover:underline text-novel">
                      {novel.source}
                    </Link>
                  </p>
                )}
              </div>
              
              {novel.genres && novel.genres.length > 0 && (
                <div>
                  <p className="text-sm text-muted-foreground mb-2">Genres:</p>
                  <div className="flex flex-wrap gap-2">
                    {novel.genres.map(genre => (
                      <Link to={`/genre/${genre}`} key={genre}>
                        <Badge variant="secondary" className="hover:bg-secondary cursor-pointer">
                          {genre}
                        </Badge>
                      </Link>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
          
          {/* Right Column - Details and Chapters */}
          <div className="md:w-2/3 lg:w-3/4">
            <h1 className="text-3xl font-bold mb-2">{novel.title}</h1>
            
            {novel.originalTitle && (
              <p className="text-lg text-muted-foreground mb-2">{novel.originalTitle}</p>
            )}
            
            {novel.author && (
              <p className="mb-4">by <span className="font-medium">{novel.author}</span></p>
            )}
            
            <p className="mb-6">{novel.summary}</p>
            
            <Separator className="my-6" />
            
            <h2 className="text-2xl font-bold mb-4">Chapters</h2>
            <ChapterList chapters={chapters} novelId={id} />
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default NovelDetail;
