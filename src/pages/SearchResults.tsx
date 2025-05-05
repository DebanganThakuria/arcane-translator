
import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import Layout from '../components/Layout';
import NovelGrid from '../components/NovelGrid';
import { searchNovels } from '../data/mockData';
import { Novel } from '../types/novel';

const SearchResults = () => {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('q') || '';
  const [results, setResults] = useState<Novel[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  
  useEffect(() => {
    // Set loading state
    setIsLoading(true);
    
    // Simulate API call with a delay
    const timer = setTimeout(() => {
      const searchResults = searchNovels(query);
      setResults(searchResults);
      setIsLoading(false);
    }, 500);
    
    return () => clearTimeout(timer);
  }, [query]);

  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <h1 className="text-3xl font-bold mb-2">Search Results</h1>
        <p className="text-muted-foreground mb-8">
          {isLoading 
            ? 'Searching...' 
            : results.length > 0 
              ? `Found ${results.length} result${results.length !== 1 ? 's' : ''} for "${query}"` 
              : `No results found for "${query}"`
          }
        </p>
        
        {isLoading ? (
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-6 animate-pulse">
            {[1, 2, 3, 4, 5, 6, 7, 8].map(i => (
              <div key={i} className="flex flex-col space-y-3">
                <div className="bg-gray-200 aspect-[2/3] rounded-md"></div>
                <div className="bg-gray-200 h-4 w-2/3 rounded"></div>
                <div className="bg-gray-200 h-3 w-1/2 rounded"></div>
              </div>
            ))}
          </div>
        ) : (
          <NovelGrid novels={results} />
        )}
      </div>
    </Layout>
  );
};

export default SearchResults;
