
import React from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { mockNovels } from '../data/mockData';

const KoreanNovels = () => {
  // Filter novels that are Korean or use mock data for demonstration
  const koreanNovels = mockNovels.filter(novel => novel.language === 'Korean' || novel.id === 'n2');

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <h1 className="text-3xl font-bold mb-8">Korean Novels</h1>
        
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Popular Korean Novels</h2>
          <NovelGrid novels={koreanNovels} />
        </div>
        
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Recently Updated Korean Novels</h2>
          <NovelGrid novels={koreanNovels.slice(0, 5)} />
        </div>
        
        <div>
          <h2 className="text-xl font-semibold mb-4">Top Korean Novel Sources</h2>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {['Munpia', 'Naver Series', 'Kakaopage', 'Ridibooks'].map(source => (
              <div key={source} className="bg-accent/50 rounded-lg p-6 text-center cursor-pointer hover:bg-accent transition-colors">
                <h3 className="font-medium">{source}</h3>
              </div>
            ))}
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default KoreanNovels;
