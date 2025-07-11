import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { searchNovelsByQuery } from '../database/db';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { Novel } from '../types/novel';

const SearchResults = () => {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('q') || '';
  const [results, setResults] = useState<Novel[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  
  useEffect(() => {
    const searchNovels = async () => {
      if (!query.trim()) {
        setResults([]);
        setIsLoading(false);
        return;
      }

      setIsLoading(true);
      
      try {
        const searchResults = await searchNovelsByQuery(query);
        setResults(searchResults);
      } catch (error) {
        console.error('Error searching novels:', error);
        setResults([]);
      } finally {
        setIsLoading(false);
      }
    };
    
    // Add a small debounce to prevent too many API calls
    const timer = setTimeout(() => {
      searchNovels();
    }, 300);
    
    return () => clearTimeout(timer);
  }, [query]);

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <h1 className="text-3xl font-bold mb-2 text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">Search Results</h1>
        <p className="text-muted-foreground mb-8">
          {isLoading 
            ? 'Searching...' 
            : results.length > 0 
              ? `Found ${results.length} result${results.length !== 1 ? 's' : ''} for "${query}"` 
              : query.trim()
                ? `No results found for "${query}"`
                : 'Enter a search term to find novels'
          }
        </p>
        
        <div className="glass-card rounded-lg p-6">
          {isLoading ? (
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-6">
              {[1, 2, 3, 4, 5, 6, 7, 8].map(i => (
                <div key={i} className="flex flex-col space-y-3">
                  <div className="bg-indigo-100 aspect-[2/3] rounded-md"></div>
                  <div className="bg-indigo-100 h-4 w-2/3 rounded"></div>
                  <div className="bg-indigo-100 h-3 w-1/2 rounded"></div>
                </div>
              ))}
            </div>
          ) : (
            <NovelGrid novels={results} />
          )}
        </div>
      </div>
    </Layout>
  );
};

export default SearchResults;
