
import React, { useEffect, useState } from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { getAllNovels } from '../database/db';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { Plus } from 'lucide-react';
import { Novel } from '../types/novel';
import { useToast } from '@/components/ui/use-toast';

const Library = () => {
  const [novels, setNovels] = useState<Novel[]>([]);
  const [loading, setLoading] = useState(true);
  const { toast } = useToast();

  useEffect(() => {
    try {
      const fetchedNovels = getAllNovels();
      setNovels(fetchedNovels);
    } catch (error) {
      console.error('Error fetching novels:', error);
      toast({
        title: "Error",
        description: "Failed to load your library.",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  }, [toast]);

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">Your Library</h1>
          <Button asChild className="bg-gradient-to-r from-indigo-600 to-blue-400 hover:from-indigo-700 hover:to-blue-500 border-none">
            <Link to="/add">
              <Plus className="mr-2 h-4 w-4" />
              Add Novel
            </Link>
          </Button>
        </div>

        {loading ? (
          <div className="glass-card rounded-lg p-10 text-center">
            <p className="text-xl text-muted-foreground mb-4">Loading your library...</p>
          </div>
        ) : novels.length > 0 ? (
          <div className="glass-card rounded-lg p-6">
            <NovelGrid novels={novels} />
          </div>
        ) : (
          <div className="glass-card rounded-lg p-10 text-center">
            <p className="text-xl text-muted-foreground mb-4">Your library is empty</p>
            <Button asChild className="bg-gradient-to-r from-indigo-600 to-blue-400 hover:from-indigo-700 hover:to-blue-500 border-none">
              <Link to="/add">Add Your First Novel</Link>
            </Button>
          </div>
        )}
      </div>
    </Layout>
  );
};

export default Library;
