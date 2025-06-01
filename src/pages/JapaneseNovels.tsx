import React, { useEffect, useState, useCallback } from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { getNovelsByFilter } from '../database/db';
import { PaginatedNovels } from '../types/novel';
import { Button } from '@/components/ui/button';
import { RefreshCw, ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-react';
import { useToast } from '@/components/ui/use-toast';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

const ITEMS_PER_PAGE_OPTIONS = [10, 20, 50, 100];

const JapaneseNovels = () => {
  const [pagination, setPagination] = useState<PaginatedNovels>({
    novels: [],
    total_count: 0,
    current_page: 1,
    total_pages: 1,
  });
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [itemsPerPage, setItemsPerPage] = useState(20);
  const [pageInput, setPageInput] = useState('1');
  const { toast } = useToast();

  const fetchNovels = useCallback(async (page = 1, limit = itemsPerPage) => {
    try {
      setLoading(true);
      const result = await getNovelsByFilter('language', 'japanese', page, limit);
      setPagination(result);
      setPageInput(String(result.current_page));
    } catch (error) {
      console.error('Error fetching novels:', error);
      toast({
        title: "Error",
        description: "Failed to load Japanese novels.",
        variant: "destructive",
      });
      setPagination(prev => ({ ...prev, novels: [] }));
    } finally {
      setLoading(false);
    }
  }, [itemsPerPage, toast]);

  // Update novels when items per page changes
  useEffect(() => {
    fetchNovels(1, itemsPerPage);
  }, [fetchNovels, itemsPerPage]);

  const handlePageChange = (newPage: number) => {
    if (newPage >= 1 && newPage <= pagination.total_pages) {
      fetchNovels(newPage, itemsPerPage);
    }
  };

  const handlePageInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    if (value === '' || /^\d+$/.test(value)) {
      setPageInput(value);
    }
  };

  const handlePageInputBlur = () => {
    const page = parseInt(pageInput, 10);
    if (!isNaN(page) && page >= 1 && page <= pagination.total_pages) {
      handlePageChange(page);
    } else {
      setPageInput(String(pagination.current_page));
    }
  };

  const handlePageInputKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handlePageInputBlur();
    }
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await fetchNovels(pagination.current_page, itemsPerPage);
      toast({
        title: "Refreshed",
        description: "Japanese novels list has been updated.",
      });
    } catch (error) {
      toast({
        title: "Refresh Failed",
        description: "Failed to refresh the novels list.",
        variant: "destructive",
      });
    } finally {
      setRefreshing(false);
    }
  };

  useEffect(() => {
    fetchNovels();
  }, [fetchNovels]);

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">
            Japanese Novels
          </h1>
          <Button 
            variant="outline" 
            onClick={handleRefresh}
            disabled={refreshing || loading}
            className="border-indigo-200 hover:bg-indigo-50"
          >
            <RefreshCw className={`mr-2 h-4 w-4 ${refreshing ? 'animate-spin' : ''}`} />
            {refreshing ? 'Updating...' : 'Refresh'}
          </Button>
        </div>
        
        <div className="space-y-6">
          {loading ? (
            <div className="glass-card rounded-lg p-10 text-center">
              <p className="text-xl text-muted-foreground mb-4">Loading Japanese novels...</p>
            </div>
          ) : pagination.novels.length > 0 ? (
            <div className="glass-card rounded-lg p-6">
              <h2 className="text-xl font-semibold mb-4">Japanese Novels</h2>
              <NovelGrid novels={pagination.novels} />
              
              {/* Pagination Controls */}
              {pagination.total_pages > 1 && (
                <div className="flex flex-col sm:flex-row items-center justify-between gap-4 mt-6">
                  <div className="flex items-center space-x-2">
                    <span className="text-sm text-muted-foreground">
                      Showing {pagination.novels.length} of {pagination.total_count} novels
                    </span>
                  </div>
                  
                  <div className="flex items-center space-x-2">
                    <Select
                      value={String(itemsPerPage)}
                      onValueChange={(value) => setItemsPerPage(Number(value))}
                    >
                      <SelectTrigger className="w-[100px] h-9">
                        <SelectValue placeholder="Items per page" />
                      </SelectTrigger>
                      <SelectContent>
                        {ITEMS_PER_PAGE_OPTIONS.map((option) => (
                          <SelectItem key={option} value={String(option)}>
                            {option} / page
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>

                    <div className="flex items-center space-x-1">
                      <Button
                        variant="outline"
                        className="h-9 w-9 p-0"
                        onClick={() => handlePageChange(1)}
                        disabled={pagination.current_page === 1 || loading}
                      >
                        <ChevronsLeft className="h-4 w-4" />
                        <span className="sr-only">First page</span>
                      </Button>
                      <Button
                        variant="outline"
                        className="h-9 w-9 p-0"
                        onClick={() => handlePageChange(pagination.current_page - 1)}
                        disabled={pagination.current_page === 1 || loading}
                      >
                        <ChevronLeft className="h-4 w-4" />
                        <span className="sr-only">Previous page</span>
                      </Button>
                      <div className="flex items-center space-x-1">
                        <span className="text-sm px-2">Page</span>
                        <Input
                          className="w-16 h-9 text-center"
                          value={pageInput}
                          onChange={handlePageInputChange}
                          onBlur={handlePageInputBlur}
                          onKeyDown={handlePageInputKeyDown}
                          disabled={loading}
                        />
                        <span className="text-sm px-2">of {pagination.total_pages}</span>
                      </div>
                      <Button
                        variant="outline"
                        className="h-9 w-9 p-0"
                        onClick={() => handlePageChange(pagination.current_page + 1)}
                        disabled={pagination.current_page === pagination.total_pages || loading}
                      >
                        <ChevronRight className="h-4 w-4" />
                        <span className="sr-only">Next page</span>
                      </Button>
                      <Button
                        variant="outline"
                        className="h-9 w-9 p-0"
                        onClick={() => handlePageChange(pagination.total_pages)}
                        disabled={pagination.current_page === pagination.total_pages || loading}
                      >
                        <ChevronsRight className="h-4 w-4" />
                        <span className="sr-only">Last page</span>
                      </Button>
                    </div>
                  </div>
                </div>
              )}
            </div>
          ) : (
            <div className="glass-card rounded-lg p-10 text-center">
              <p className="text-xl text-muted-foreground mb-4">No Japanese novels found</p>
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default JapaneseNovels;
