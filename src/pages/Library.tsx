
import React from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { mockNovels } from '../data/mockData';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { Plus } from 'lucide-react';

const Library = () => {
  // In a real app, this would fetch from an API
  const novels = mockNovels;

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

        {novels.length > 0 ? (
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
