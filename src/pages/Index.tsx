
import React from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { getRecentNovels, getRecentlyUpdatedNovels, mockNovels } from '../data/mockData';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

const Index = () => {
  const recentNovels = getRecentNovels();
  const updatedNovels = getRecentlyUpdatedNovels();
  
  // For demonstration, we'll filter novels by language
  // In a real app, this would come from the API
  const chineseNovels = mockNovels.filter(novel => novel.language === 'Chinese' || novel.id === 'n1');
  const koreanNovels = mockNovels.filter(novel => novel.language === 'Korean' || novel.id === 'n2');
  const japaneseNovels = mockNovels.filter(novel => novel.language === 'Japanese' || novel.id === 'n3');

  return (
    <Layout>
      <section className="py-10 px-4">
        <div className="container mx-auto">
          {/* Hero Section */}
          <div className="relative overflow-hidden rounded-xl mb-12">
            <div className="bg-novel-dark/80 text-white p-8 sm:p-12 rounded-xl">
              <div className="max-w-3xl mx-auto">
                <h1 className="text-3xl sm:text-4xl md:text-5xl font-bold mb-4">
                  Your Favorite Webnovels, Translated With AI
                </h1>
                <p className="text-lg mb-6 opacity-90">
                  Read Chinese, Korean, and Japanese webnovels seamlessly translated 
                  into English, right as you read.
                </p>
                <div className="flex flex-wrap gap-4">
                  <Button asChild size="lg" className="bg-novel hover:bg-novel/90">
                    <Link to="/add">Add Your First Novel</Link>
                  </Button>
                  <Button asChild size="lg" variant="outline" className="border-white/20 bg-white/10 hover:bg-white/20">
                    <Link to="/library">Browse Library</Link>
                  </Button>
                </div>
              </div>
            </div>
          </div>

          {/* Continue Reading Section */}
          {recentNovels.length > 0 && (
            <div className="mb-12">
              <div className="flex justify-between items-center mb-6">
                <h2 className="text-2xl font-bold">Continue Reading</h2>
                <Button asChild variant="ghost">
                  <Link to="/library">View All</Link>
                </Button>
              </div>
              <NovelGrid novels={recentNovels} recent={true} />
            </div>
          )}

          {/* Novels by Language */}
          <div className="mb-12">
            <h2 className="text-2xl font-bold mb-6">Browse by Language</h2>
            <Tabs defaultValue="chinese" className="w-full">
              <TabsList className="mb-6">
                <TabsTrigger value="chinese">Chinese</TabsTrigger>
                <TabsTrigger value="korean">Korean</TabsTrigger>
                <TabsTrigger value="japanese">Japanese</TabsTrigger>
              </TabsList>
              
              <TabsContent value="chinese">
                <div className="mb-4">
                  <h3 className="text-lg font-medium mb-4">Popular Chinese Novels</h3>
                  <NovelGrid novels={chineseNovels} />
                </div>
              </TabsContent>
              
              <TabsContent value="korean">
                <div className="mb-4">
                  <h3 className="text-lg font-medium mb-4">Popular Korean Novels</h3>
                  <NovelGrid novels={koreanNovels} />
                </div>
              </TabsContent>
              
              <TabsContent value="japanese">
                <div className="mb-4">
                  <h3 className="text-lg font-medium mb-4">Popular Japanese Novels</h3>
                  <NovelGrid novels={japaneseNovels} />
                </div>
              </TabsContent>
            </Tabs>
          </div>

          {/* Recently Updated Section */}
          <div className="mb-12">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-bold">Recently Updated</h2>
              <Button asChild variant="ghost">
                <Link to="/library">View All</Link>
              </Button>
            </div>
            <NovelGrid novels={updatedNovels} />
          </div>
        </div>
      </section>
    </Layout>
  );
};

export default Index;
