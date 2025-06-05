import React, { useEffect, useCallback } from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import PaginationControls from '../components/PaginationControls';
import { getAllNovels } from '../database/db';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { Plus } from 'lucide-react';
import { usePagination } from '../hooks/usePagination';

const Library = () => {
  const fetchNovels = useCallback((page: number, limit: number) => {
    return getAllNovels(page, limit);
  }, []);

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
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">Your Library</h1>
          
          <div className="flex space-x-2">
            <Button 
              variant="outline" 
              onClick={handleRefresh}
              disabled={refreshing || loading}
              className="border-indigo-200 hover:bg-indigo-50"
            >
              {refreshing ? 'Updating...' : 'Refresh'}
            </Button>
            
            <Button asChild className="bg-gradient-to-r from-indigo-600 to-blue-400 hover:from-indigo-700 hover:to-blue-500 border-none">
              <Link to="/add">
                <Plus className="mr-2 h-4 w-4" />
                Add Novel
              </Link>
            </Button>
          </div>
        </div>

        <div className="space-y-6">
          {loading ? (
            <div className="glass-card rounded-lg p-10 text-center">
              <p className="text-xl text-muted-foreground mb-4">Loading your library...</p>
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
              <p className="text-xl text-muted-foreground mb-4">Your library is empty</p>
              <Button asChild className="bg-gradient-to-r from-indigo-600 to-blue-400 hover:from-indigo-700 hover:to-blue-500 border-none">
                <Link to="/add">Add Your First Novel</Link>
              </Button>
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default Library;
