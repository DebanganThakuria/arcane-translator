
import React, { useEffect, useState } from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { getAllNovels } from '../database/db';
import { Novel } from '../types/novel';
import { Link } from 'react-router-dom';

const KoreanNovels = () => {
  const [novels, setNovels] = useState<Novel[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchNovels = async () => {
      try {
        const fetchedNovels = await getAllNovels();
        // Filter novels that are from Korean sources
        const koreanNovels = fetchedNovels.filter(novel => 
          novel.source === 'naver' || novel.source === 'munpia' || novel.source === 'kakaopage'
        );
        setNovels(koreanNovels);
      } catch (error) {
        console.error('Error fetching novels:', error);
        setNovels([]);
      } finally {
        setLoading(false);
      }
    };

    fetchNovels();
  }, []);

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
        <h1 className="text-3xl font-bold mb-8 text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">
          Korean Novels
        </h1>
        
        <div className="glass-card rounded-lg p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Popular Korean Novels</h2>
          <NovelGrid novels={novels} />
        </div>
        
        <div className="glass-card rounded-lg p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Recently Updated Korean Novels</h2>
          <NovelGrid novels={novels.slice(0, 3)} />
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
