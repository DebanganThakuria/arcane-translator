
import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { getAllNovels } from '../database/db';
import { Novel } from '../types/novel';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';

const GenrePage = () => {
  const { genre } = useParams<{ genre: string }>();
  const [novels, setNovels] = useState<Novel[]>([]);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    const fetchNovels = async () => {
      try {
        const fetchedNovels = await getAllNovels();
        // Filter novels by genre
        const filteredNovels = fetchedNovels.filter(novel => 
          novel.genres && novel.genres.includes(genre || '')
        );
        setNovels(filteredNovels);
      } catch (error) {
        console.error('Error fetching novels:', error);
        setNovels([]);
      } finally {
        setLoading(false);
      }
    };

    fetchNovels();
  }, [genre]);

  if (loading) {
    return (
      <Layout>
        <div className="container mx-auto py-16 px-4 text-center">
          <p className="text-xl text-muted-foreground mb-4">Loading...</p>
        </div>
      </Layout>
    );
  }
  
  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <Button asChild variant="ghost" className="mb-6 text-indigo-600 hover:text-indigo-800 hover:bg-indigo-50">
          <Link to="/library">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Library
          </Link>
        </Button>
      
        <h1 className="text-3xl font-bold mb-8 text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">
          {genre || 'Unknown Genre'} Novels
        </h1>
        
        <div className="glass-card rounded-lg p-6">
          {novels.length > 0 ? (
            <NovelGrid novels={novels} />
          ) : (
            <div className="text-center py-12">
              <p className="text-muted-foreground">No novels found in this genre.</p>
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default GenrePage;
