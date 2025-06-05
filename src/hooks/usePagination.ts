
import { useState, useCallback } from 'react';
import { PaginatedNovels } from '../types/novel';
import { useToast } from '@/components/ui/use-toast';

interface UsePaginationProps {
  fetchFunction: (page: number, limit: number) => Promise<PaginatedNovels>;
  itemsPerPage?: number;
}

export const usePagination = ({ fetchFunction, itemsPerPage = 20 }: UsePaginationProps) => {
  const [pagination, setPagination] = useState<PaginatedNovels>({
    novels: [],
    total_count: 0,
    current_page: 1,
    total_pages: 1,
  });
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [itemsPerPageState, setItemsPerPageState] = useState(itemsPerPage);
  const [pageInput, setPageInput] = useState('1');
  const { toast } = useToast();

  const fetchData = useCallback(async (page = 1, limit = itemsPerPageState) => {
    try {
      setLoading(true);
      const result = await fetchFunction(page, limit);
      setPagination(result);
      setPageInput(String(result.current_page));
    } catch (error) {
      console.error('Error fetching data:', error);
      toast({
        title: "Error",
        description: "Failed to load data.",
        variant: "destructive",
      });
      setPagination(prev => ({ ...prev, novels: [] }));
    } finally {
      setLoading(false);
    }
  }, [fetchFunction, itemsPerPageState, toast]);

  const handlePageChange = (newPage: number) => {
    if (newPage >= 1 && newPage <= pagination.total_pages) {
      fetchData(newPage, itemsPerPageState);
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
      await fetchData(pagination.current_page, itemsPerPageState);
      toast({
        title: "Refreshed",
        description: "Data has been updated.",
      });
    } catch (error) {
      toast({
        title: "Refresh Failed",
        description: "Failed to refresh data.",
        variant: "destructive",
      });
    } finally {
      setRefreshing(false);
    }
  };

  const setItemsPerPage = (newItemsPerPage: number) => {
    setItemsPerPageState(newItemsPerPage);
    fetchData(1, newItemsPerPage);
  };

  return {
    pagination,
    loading,
    refreshing,
    itemsPerPage: itemsPerPageState,
    pageInput,
    fetchData,
    handlePageChange,
    handlePageInputChange,
    handlePageInputBlur,
    handlePageInputKeyDown,
    handleRefresh,
    setItemsPerPage,
  };
};
