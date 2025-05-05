
import React from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { getRecentNovels, getRecentlyUpdatedNovels, mockNovels } from '../data/mockData';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';

const Index = () => {
  const recentNovels = getRecentNovels();
  const updatedNovels = getRecentlyUpdatedNovels();
  const featuredNovel = mockNovels[2]; // Using Silent Shadows as featured

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

          {/* Featured Novel */}
          {featuredNovel && (
            <div className="mb-12 bg-gradient-to-r from-novel-light/20 to-novel/10 p-6 rounded-xl">
              <div className="flex flex-col md:flex-row gap-8">
                <div className="md:w-1/3">
                  <div className="book-cover aspect-[2/3] max-w-[250px] mx-auto">
                    <img 
                      src={featuredNovel.cover || '/placeholder.svg'} 
                      alt={featuredNovel.title} 
                      className="w-full h-full object-cover rounded-md shadow-lg"
                    />
                  </div>
                </div>
                <div className="md:w-2/3">
                  <h3 className="text-sm font-medium text-novel mb-2">FEATURED NOVEL</h3>
                  <h2 className="text-3xl font-bold mb-4">{featuredNovel.title}</h2>
                  <p className="text-sm mb-2">
                    <span className="text-muted-foreground">Original Title: </span>
                    {featuredNovel.originalTitle}
                  </p>
                  <p className="text-sm mb-4">
                    <span className="text-muted-foreground">Author: </span>
                    {featuredNovel.author}
                  </p>
                  <p className="mb-6">{featuredNovel.summary}</p>
                  <Button asChild>
                    <Link to={`/novel/${featuredNovel.id}`}>View Details</Link>
                  </Button>
                </div>
              </div>
            </div>
          )}

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
