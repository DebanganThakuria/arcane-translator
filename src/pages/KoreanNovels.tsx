
import React from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { mockNovels } from '../data/mockData';
import { Link } from 'react-router-dom';

const KoreanNovels = () => {
  // Filter novels for Korean - since language property isn't available, 
  // we're using a mock approach for demonstration
  const koreanNovels = mockNovels.filter(novel => 
    // Use ID to filter for demo purposes since language isn't available in the Novel type
    novel.id === 'n2' || novel.id === 'n5'
  );

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <h1 className="text-3xl font-bold mb-8 text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">
          Korean Novels
        </h1>
        
        <div className="glass-card rounded-lg p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Popular Korean Novels</h2>
          <NovelGrid novels={koreanNovels} />
        </div>
        
        <div className="glass-card rounded-lg p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Recently Updated Korean Novels</h2>
          <NovelGrid novels={koreanNovels.slice(0, 3)} />
        </div>
        
        <div className="glass-card rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4">Top Korean Novel Sources</h2>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {['Munpia', 'Naver Series', 'Kakaopage', 'Ridibooks'].map(source => (
              <Link key={source} to={`/source/${source}`} className="bg-white/60 hover:bg-white/80 rounded-lg p-6 text-center cursor-pointer transition-colors shadow-sm hover:shadow-md">
                <h3 className="font-medium">{source}</h3>
              </Link>
            ))}
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default KoreanNovels;
