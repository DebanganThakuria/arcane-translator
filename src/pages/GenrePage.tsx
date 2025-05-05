
import React from 'react';
import { useParams, Link } from 'react-router-dom';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { mockNovels } from '../data/mockData';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';

const GenrePage = () => {
  const { genre } = useParams<{ genre: string }>();
  
  // Filter novels by genre - in a real app, this would call an API
  const filteredNovels = mockNovels.filter(novel => 
    novel.genres && novel.genres.includes(genre || '')
  );
  
  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <Button asChild variant="ghost" className="mb-6">
          <Link to="/library">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Library
          </Link>
        </Button>
      
        <h1 className="text-3xl font-bold mb-8">
          {genre} Novels
        </h1>
        
        {filteredNovels.length > 0 ? (
          <NovelGrid novels={filteredNovels} />
        ) : (
          <div className="text-center py-12">
            <p className="text-muted-foreground">No novels found in this genre.</p>
          </div>
        )}
      </div>
    </Layout>
  );
};

export default GenrePage;
