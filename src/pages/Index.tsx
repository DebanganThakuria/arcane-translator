import React, { useEffect, useState } from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import LoadingSkeleton from '../components/LoadingSkeleton';
import { getNovelsByFilter, getNovelsStats } from '../database/db';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Novel } from '../types/novel';
import { BookOpen, TrendingUp, Clock, Globe, Sparkles, ArrowRight } from 'lucide-react';

const Index = () => {
  const [chineseNovels, setChineseNovels] = useState<Novel[]>([]);
  const [koreanNovels, setKoreanNovels] = useState<Novel[]>([]);
  const [japaneseNovels, setJapaneseNovels] = useState<Novel[]>([]);
  const [recentlyReadNovels, setRecentlyReadNovels] = useState<Novel[]>([]);
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    novel_count: 0,
    chapter_count: 0,
    languages_supported: 3,
  });

  useEffect(() => {
    const fetchNovels = async () => {
      try {
        setLoading(true);

        const [chinese, korean, japanese, recentlyRead] = await Promise.all([
          getNovelsByFilter('language', 'chinese', 1, 5),
          getNovelsByFilter('language', 'korean', 1, 5),
          getNovelsByFilter('language', 'japanese', 1, 5),
          getNovelsByFilter('recently_read', '5')
        ]);

        setChineseNovels(chinese.novels);
        setKoreanNovels(korean.novels);
        setJapaneseNovels(japanese.novels);
        setRecentlyReadNovels(recentlyRead.novels);

        // Calculate stats
        const stats = await getNovelsStats();

        console.log(stats)

        setStats({
          novel_count: stats.novel_count,
          chapter_count: stats.chapter_count,
          languages_supported: 3
        });
      } catch (error) {
        console.error('Error fetching novels:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchNovels();
  }, []);

  const StatCard = ({ icon: Icon, title, value, description }: {
    icon: React.ElementType;
    title: string;
    value: string | number;
    description: string;
  }) => (
    <Card className="stat-card">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium text-gray-700 dark:text-gray-200">{title}</CardTitle>
        <Icon className="h-4 w-4 text-indigo-500 dark:text-indigo-400" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold gradient-text">{value}</div>
        <p className="text-xs text-gray-600 dark:text-gray-300">{description}</p>
      </CardContent>
    </Card>
  );

  if (loading) {
    return (
      <Layout>
        <div className="container mx-auto py-16 px-4">
          {/* Hero Skeleton */}
          <div className="glass-card rounded-xl p-8 sm:p-12 mb-12">
            <div className="max-w-3xl mx-auto space-y-4">
              <LoadingSkeleton variant="text" count={3} />
              <div className="flex gap-4 mt-6">
                <div className="skeleton h-10 w-32 rounded-lg" />
                <div className="skeleton h-10 w-32 rounded-lg" />
              </div>
            </div>
          </div>

          {/* Stats Skeleton */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-12">
            {Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="glass-card rounded-lg p-4">
                <LoadingSkeleton variant="text" count={2} />
              </div>
            ))}
          </div>

          {/* Content Skeleton */}
          <div className="space-y-12">
            <div className="glass-card rounded-lg p-6">
              <div className="skeleton h-6 w-48 mb-6 rounded" />
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
                <LoadingSkeleton variant="card" count={6} />
              </div>
            </div>
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <section className="py-10 px-4">
        <div className="container mx-auto">
          {/* Hero Section */}
          <div className="relative overflow-hidden rounded-xl mb-12">
            <div className="glass-card bg-gradient-to-br from-indigo-600/90 to-blue-500/90 text-white p-8 sm:p-12 rounded-xl">
              <div className="max-w-3xl mx-auto text-center">
                <div className="inline-block text-4xl mb-4">✨</div>
                <h1 className="text-3xl sm:text-4xl md:text-5xl font-bold mb-4">
                  Your Favorite Webnovels, <br />
                  <span className="text-yellow-300">Translated With AI</span>
                </h1>
                <p className="text-lg mb-8 opacity-90">
                  Read Chinese, Korean, and Japanese webnovels seamlessly translated 
                  into English, right as you read. Experience stories like never before.
                </p>
                <div className="flex flex-wrap gap-4 justify-center">
                  <Button asChild size="lg" className="bg-white text-indigo-600 hover:bg-gray-100 font-semibold">
                    <Link to="/add">
                      <Sparkles className="mr-2 h-5 w-5" />
                      Add Your First Novel
                    </Link>
                  </Button>
                  <Button asChild size="lg" variant="outline" className="border-white/30 bg-white/10 hover:bg-white/20 text-white">
                    <Link to="/library">
                      <BookOpen className="mr-2 h-5 w-5" />
                      Browse Library
                    </Link>
                  </Button>
                </div>
              </div>
            </div>
          </div>

          {/* Statistics Section */}
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mb-12">
            <StatCard
              icon={BookOpen}
              title="Total Novels"
              value={stats.novel_count}
              description="In your library"
            />
            <StatCard
              icon={TrendingUp}
              title="Chapters"
              value={stats.chapter_count}
              description="Ready to read"
            />
            <StatCard
              icon={Globe}
              title="Languages"
              value={stats.languages_supported}
              description="Supported"
            />
          </div>

          {/* Recently Read Section */}
          <div className="mb-12">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-bold gradient-text">Recently Read</h2>
              <Button asChild variant="ghost" className="text-indigo-600 hover:text-indigo-800 hover:bg-indigo-50 dark:text-indigo-400 dark:hover:text-indigo-300">
                <Link to="/library">
                  View All
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
            <div className="glass-card rounded-lg p-6">
              {recentlyReadNovels.length > 0 ? (
                <NovelGrid novels={recentlyReadNovels} />
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  <Clock className="h-12 w-12 mx-auto mb-4 opacity-50" />
                  <p>No recently read novels. <Link to="/add" className="text-indigo-600 hover:underline">Add some novels to get started!</Link></p>
                </div>
              )}
            </div>
          </div>
          
          {/* Browse by Language Section */}
          <div className="mb-12">
            <h2 className="text-2xl font-bold mb-6 gradient-text">Browse by Language</h2>
            <div className="glass-card rounded-lg p-6">
              <Tabs defaultValue="chinese" className="w-full">
                <TabsList className="mb-6 bg-white/80 backdrop-blur-sm dark:bg-gray-700/80 border border-gray-200 dark:border-gray-600">
                  <TabsTrigger value="chinese" className="data-[state=active]:bg-indigo-100 dark:data-[state=active]:bg-indigo-900 data-[state=active]:text-indigo-700 dark:data-[state=active]:text-indigo-300">
                    <span className="mr-2">🇨🇳</span>
                    Chinese
                  </TabsTrigger>
                  <TabsTrigger value="korean" className="data-[state=active]:bg-indigo-100 dark:data-[state=active]:bg-indigo-900 data-[state=active]:text-indigo-700 dark:data-[state=active]:text-indigo-300">
                    <span className="mr-2">🇰🇷</span>
                    Korean
                  </TabsTrigger>
                  <TabsTrigger value="japanese" className="data-[state=active]:bg-indigo-100 dark:data-[state=active]:bg-indigo-900 data-[state=active]:text-indigo-700 dark:data-[state=active]:text-indigo-300">
                    <span className="mr-2">🇯🇵</span>
                    Japanese
                  </TabsTrigger>
                </TabsList>
                
                <TabsContent value="chinese">
                  <div className="mb-4">
                    <div className="flex justify-between items-center mb-4">
                      <h3 className="text-lg font-medium">Chinese Novels</h3>
                      <Button asChild variant="outline" size="sm">
                        <Link to="/chinese">
                          View All
                          <ArrowRight className="ml-2 h-4 w-4" />
                        </Link>
                      </Button>
                    </div>
                    {chineseNovels.length > 0 ? (
                      <NovelGrid novels={chineseNovels} />
                    ) : (
                      <div className="text-center py-8 text-muted-foreground">
                        <BookOpen className="h-12 w-12 mx-auto mb-4 opacity-50" />
                        <p>No Chinese novels yet. <Link to="/add" className="text-indigo-600 hover:underline">Add one now!</Link></p>
                      </div>
                    )}
                  </div>
                </TabsContent>
                
                <TabsContent value="korean">
                  <div className="mb-4">
                    <div className="flex justify-between items-center mb-4">
                      <h3 className="text-lg font-medium">Korean Novels</h3>
                      <Button asChild variant="outline" size="sm">
                        <Link to="/korean">
                          View All
                          <ArrowRight className="ml-2 h-4 w-4" />
                        </Link>
                      </Button>
                    </div>
                    {koreanNovels.length > 0 ? (
                      <NovelGrid novels={koreanNovels} />
                    ) : (
                      <div className="text-center py-8 text-muted-foreground">
                        <BookOpen className="h-12 w-12 mx-auto mb-4 opacity-50" />
                        <p>No Korean novels yet. <Link to="/add" className="text-indigo-600 hover:underline">Add one now!</Link></p>
                      </div>
                    )}
                  </div>
                </TabsContent>
                
                <TabsContent value="japanese">
                  <div className="mb-4">
                    <div className="flex justify-between items-center mb-4">
                      <h3 className="text-lg font-medium">Japanese Novels</h3>
                      <Button asChild variant="outline" size="sm">
                        <Link to="/japanese">
                          View All
                          <ArrowRight className="ml-2 h-4 w-4" />
                        </Link>
                      </Button>
                    </div>
                    {japaneseNovels.length > 0 ? (
                      <NovelGrid novels={japaneseNovels} />
                    ) : (
                      <div className="text-center py-8 text-muted-foreground">
                        <BookOpen className="h-12 w-12 mx-auto mb-4 opacity-50" />
                        <p>No Japanese novels yet. <Link to="/add" className="text-indigo-600 hover:underline">Add one now!</Link></p>
                      </div>
                    )}
                  </div>
                </TabsContent>
              </Tabs>
            </div>
          </div>

        </div>
      </section>
    </Layout>
  );
};

export default Index;
