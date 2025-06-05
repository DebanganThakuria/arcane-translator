import React, { useEffect, useCallback } from 'react';
import { useParams, Link } from 'react-router-dom';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import PaginationControls from '../components/PaginationControls';
import { getNovelsByFilter } from '../database/db';
import { Button } from '@/components/ui/button';
import { ArrowLeft, RefreshCw } from 'lucide-react';
import { usePagination } from '../hooks/usePagination';

const GenrePage = () => {
  const { genre } = useParams<{ genre: string }>();

  const fetchNovels = useCallback((page: number, limit: number) => {
    return getNovelsByFilter('genre', genre || '', page, limit);
  }, [genre]);

  const {
    pagination,
    loading,
    refreshing,
    itemsPerPage,
    pageInput,
    fetchData,
    handlePageChange,
    handlePageInputChange,
    handlePageInputBlur,
    handlePageInputKeyDown,
    handleRefresh,
    setItemsPerPage,
  } = usePagination({ fetchFunction: fetchNovels });

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-6 gap-4">
          <Button asChild variant="ghost" className="text-indigo-600 hover:text-indigo-800 hover:bg-indigo-50">
            <Link to="/library">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Library
            </Link>
          </Button>
          
          <div className="flex items-center gap-2">
            <h1 className="text-2xl sm:text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">
              {genre || 'Unknown Genre'} Novels
            </h1>
            <Button 
              variant="outline" 
              size="icon"
              onClick={handleRefresh}
              disabled={refreshing || loading}
              className="border-indigo-200 hover:bg-indigo-50"
            >
              <RefreshCw className={`h-4 w-4 ${refreshing ? 'animate-spin' : ''}`} />
              <span className="sr-only">Refresh</span>
            </Button>
          </div>
        </div>
        
        <div className="space-y-6">
          {loading ? (
            <div className="glass-card rounded-lg p-10 text-center">
              <p className="text-xl text-muted-foreground mb-4">
                Loading {genre ? `${genre} ` : ''}novels...
              </p>
            </div>
          ) : pagination.novels.length > 0 ? (
            <div className="glass-card rounded-lg p-6">
              <NovelGrid novels={pagination.novels} />
              <PaginationControls
                pagination={pagination}
                loading={loading}
                itemsPerPage={itemsPerPage}
                pageInput={pageInput}
                onPageChange={handlePageChange}
                onPageInputChange={handlePageInputChange}
                onPageInputBlur={handlePageInputBlur}
                onPageInputKeyDown={handlePageInputKeyDown}
                onItemsPerPageChange={setItemsPerPage}
              />
            </div>
          ) : (
            <div className="glass-card rounded-lg p-10 text-center">
              <p className="text-xl text-muted-foreground mb-4">
                No {genre ? `${genre} ` : ''}novels found.
              </p>
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default GenrePage;
