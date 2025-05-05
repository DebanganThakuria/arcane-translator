
import React from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { mockNovels } from '../data/mockData';

const ChineseNovels = () => {
  // Filter novels that are Chinese or use mock data for demonstration
  const chineseNovels = mockNovels.filter(novel => novel.language === 'Chinese' || novel.id === 'n1');

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <h1 className="text-3xl font-bold mb-8">Chinese Novels</h1>
        
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Popular Chinese Novels</h2>
          <NovelGrid novels={chineseNovels} />
        </div>
        
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Recently Updated Chinese Novels</h2>
          <NovelGrid novels={chineseNovels.slice(0, 5)} />
        </div>
        
        <div>
          <h2 className="text-xl font-semibold mb-4">Top Chinese Novel Sources</h2>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {['Qidian', 'Zongheng', 'Jinjiang', 'Ciweimao'].map(source => (
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

export default ChineseNovels;
