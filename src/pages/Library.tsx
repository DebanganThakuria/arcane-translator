import React, { useEffect, useCallback, useState } from 'react';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import LoadingSkeleton from '../components/LoadingSkeleton';
import PaginationControls from '../components/PaginationControls';
import { getAllNovels } from '../database/db';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Link } from 'react-router-dom';
import { Plus, Filter, Grid, List, SortAsc, SortDesc, Search, RefreshCw } from 'lucide-react';
import { usePagination } from '../hooks/usePagination';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Badge } from '@/components/ui/badge';
import { Novel } from '../types/novel';

const Library = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<'title' | 'date_added' | 'last_updated' | 'chapters_count'>('last_updated');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
  const [filterStatus, setFilterStatus] = useState<'all' | 'ongoing' | 'completed'>('all');
  const [filterLanguage, setFilterLanguage] = useState<'all' | 'chinese' | 'korean' | 'japanese'>('all');
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [filteredNovels, setFilteredNovels] = useState<Novel[]>([]);

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

  // Filter and sort novels
  useEffect(() => {
    let filtered = [...pagination.novels];

    // Apply search filter
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(novel => 
        novel.title.toLowerCase().includes(query) ||
        novel.author?.toLowerCase().includes(query) ||
        novel.genres?.some(genre => genre.toLowerCase().includes(query))
      );
    }

    // Apply status filter
    if (filterStatus !== 'all') {
      filtered = filtered.filter(novel => 
        novel.status?.toLowerCase() === filterStatus
      );
    }

    // Apply language filter
    if (filterLanguage !== 'all') {
      filtered = filtered.filter(novel => 
        novel.source.toLowerCase().includes(filterLanguage)
      );
    }

    // Apply sorting
    filtered.sort((a, b) => {
      let aValue: any, bValue: any;
      
      switch (sortBy) {
        case 'title':
          aValue = a.title.toLowerCase();
          bValue = b.title.toLowerCase();
          break;
        case 'date_added':
          aValue = a.date_added;
          bValue = b.date_added;
          break;
        case 'last_updated':
          aValue = a.last_updated;
          bValue = b.last_updated;
          break;
        case 'chapters_count':
          aValue = a.chapters_count;
          bValue = b.chapters_count;
          break;
        default:
          return 0;
      }

      if (sortOrder === 'asc') {
        return aValue > bValue ? 1 : -1;
      } else {
        return aValue < bValue ? 1 : -1;
      }
    });

    setFilteredNovels(filtered);
  }, [pagination.novels, searchQuery, sortBy, sortOrder, filterStatus, filterLanguage]);

  const getActiveFiltersCount = () => {
    let count = 0;
    if (searchQuery.trim()) count++;
    if (filterStatus !== 'all') count++;
    if (filterLanguage !== 'all') count++;
    return count;
  };

  const clearAllFilters = () => {
    setSearchQuery('');
    setFilterStatus('all');
    setFilterLanguage('all');
    setSortBy('last_updated');
    setSortOrder('desc');
  };

  if (loading) {
    return (
      <Layout>
        <div className="container mx-auto py-8 px-4">
          <div className="flex justify-between items-center mb-8">
            <div className="skeleton h-8 w-48 rounded" />
            <div className="flex space-x-2">
              <div className="skeleton h-10 w-24 rounded" />
              <div className="skeleton h-10 w-32 rounded" />
            </div>
          </div>

          <div className="glass-card rounded-lg p-6">
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
              <LoadingSkeleton variant="card" count={12} />
            </div>
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-8 gap-4">
          <div>
            <h1 className="text-3xl font-bold gradient-text">Your Library</h1>
            <p className="text-muted-foreground mt-1">
              {pagination.total_count} novels total
            </p>
          </div>
          
          <div className="flex flex-wrap gap-2">
            <Button 
              variant="outline" 
              onClick={handleRefresh}
              disabled={refreshing || loading}
              className="border-indigo-200 hover:bg-indigo-50 dark:border-indigo-800 dark:hover:bg-indigo-950"
            >
              <RefreshCw className={`mr-2 h-4 w-4 ${refreshing ? 'animate-spin' : ''}`} />
              {refreshing ? 'Updating...' : 'Refresh'}
            </Button>
            
            <Button asChild className="btn-primary">
              <Link to="/add">
                <Plus className="mr-2 h-4 w-4" />
                Add Novel
              </Link>
            </Button>
          </div>
        </div>

        {/* Filters and Controls */}
        <div className="glass-card rounded-lg p-4 mb-6">
          <div className="flex flex-col lg:flex-row gap-4">
            {/* Search */}
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search novels, authors, genres..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>

            {/* Filters */}
            <div className="flex flex-wrap gap-2">
              <Select value={filterLanguage} onValueChange={(value: any) => setFilterLanguage(value)}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="Language" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Languages</SelectItem>
                  <SelectItem value="chinese">🇨🇳 Chinese</SelectItem>
                  <SelectItem value="korean">🇰🇷 Korean</SelectItem>
                  <SelectItem value="japanese">🇯🇵 Japanese</SelectItem>
                </SelectContent>
              </Select>

              <Select value={filterStatus} onValueChange={(value: any) => setFilterStatus(value)}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="ongoing">Ongoing</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                </SelectContent>
              </Select>

              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" className="gap-2">
                    {sortOrder === 'asc' ? <SortAsc className="h-4 w-4" /> : <SortDesc className="h-4 w-4" />}
                    Sort
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuLabel>Sort by</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => setSortBy('last_updated')}>
                    Last Updated
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setSortBy('date_added')}>
                    Date Added
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setSortBy('title')}>
                    Title
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setSortBy('chapters_count')}>
                    Chapter Count
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')}>
                    {sortOrder === 'asc' ? 'Descending' : 'Ascending'}
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>

              <div className="flex border rounded-lg">
                <Button
                  variant={viewMode === 'grid' ? 'default' : 'ghost'}
                  size="sm"
                  onClick={() => setViewMode('grid')}
                  className="rounded-r-none"
                >
                  <Grid className="h-4 w-4" />
                </Button>
                <Button
                  variant={viewMode === 'list' ? 'default' : 'ghost'}
                  size="sm"
                  onClick={() => setViewMode('list')}
                  className="rounded-l-none"
                >
                  <List className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>

          {/* Active Filters */}
          {getActiveFiltersCount() > 0 && (
            <div className="flex items-center gap-2 mt-4 pt-4 border-t">
              <span className="text-sm text-muted-foreground">Active filters:</span>
              {searchQuery.trim() && (
                <Badge variant="secondary">
                  Search: "{searchQuery}"
                </Badge>
              )}
              {filterLanguage !== 'all' && (
                <Badge variant="secondary">
                  Language: {filterLanguage}
                </Badge>
              )}
              {filterStatus !== 'all' && (
                <Badge variant="secondary">
                  Status: {filterStatus}
                </Badge>
              )}
              <Button
                variant="ghost"
                size="sm"
                onClick={clearAllFilters}
                className="text-xs"
              >
                Clear all
              </Button>
            </div>
          )}
        </div>

        {/* Results */}
        <div className="space-y-6">
          {filteredNovels.length > 0 ? (
            <div className="glass-card rounded-lg p-6">
              <div className="flex justify-between items-center mb-4">
                <p className="text-sm text-muted-foreground">
                  Showing {filteredNovels.length} of {pagination.total_count} novels
                </p>
              </div>
              
              <NovelGrid novels={filteredNovels} showStats={true} />
              
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
              {getActiveFiltersCount() > 0 ? (
                <>
                  <Filter className="h-12 w-12 mx-auto mb-4 text-muted-foreground opacity-50" />
                  <p className="text-xl text-muted-foreground mb-4">No novels match your filters</p>
                  <Button onClick={clearAllFilters} variant="outline">
                    Clear Filters
                  </Button>
                </>
              ) : (
                <>
                  <Plus className="h-12 w-12 mx-auto mb-4 text-muted-foreground opacity-50" />
                  <p className="text-xl text-muted-foreground mb-4">Your library is empty</p>
                  <Button asChild className="btn-primary">
                    <Link to="/add">Add Your First Novel</Link>
                  </Button>
                </>
              )}
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default Library;
