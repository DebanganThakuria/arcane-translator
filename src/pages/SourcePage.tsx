
import React from 'react';
import { useParams, Link } from 'react-router-dom';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { mockNovels } from '../data/mockData';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';

const SourcePage = () => {
  const { source } = useParams<{ source: string }>();
  
  // Filter novels by source - in a real app, this would call an API
  const filteredNovels = mockNovels.filter(novel => 
    novel.source === source
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
      
        <h1 className="text-3xl font-bold mb-8 text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">
          Novels from {source}
        </h1>
        
        <div className="glass-card rounded-lg p-6">
          {filteredNovels.length > 0 ? (
            <NovelGrid novels={filteredNovels} />
          ) : (
            <div className="text-center py-12">
              <p className="text-muted-foreground">No novels found from this source.</p>
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default SourcePage;
