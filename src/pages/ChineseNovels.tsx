
import React from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { mockNovels } from '../data/mockData';
import { Link } from 'react-router-dom';

const ChineseNovels = () => {
  // Filter novels that are Chinese or use mock data for demonstration
  // Since language property doesn't exist in the Novel type, we're filtering specific novels
  const chineseNovels = mockNovels.filter(novel => novel.id === 'n1' || novel.id === 'n4');

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <h1 className="text-3xl font-bold mb-8 text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">
          Chinese Novels
        </h1>
        
        <div className="glass-card rounded-lg p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Popular Chinese Novels</h2>
          <NovelGrid novels={chineseNovels} />
        </div>
        
        <div className="glass-card rounded-lg p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Recently Updated Chinese Novels</h2>
          <NovelGrid novels={chineseNovels.slice(0, 5)} />
        </div>
        
        <div className="glass-card rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4">Top Chinese Novel Sources</h2>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {['Qidian', 'Zongheng', 'Jinjiang', 'Ciweimao'].map(source => (
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

export default ChineseNovels;
