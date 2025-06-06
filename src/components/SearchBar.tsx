import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Search, Clock, X } from 'lucide-react';
import {
  Command,
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command';

const SearchBar = () => {
  const [query, setQuery] = useState('');
  const [open, setOpen] = useState(false);
  const [recentSearches, setRecentSearches] = useState<string[]>([]);
  const navigate = useNavigate();
  const inputRef = useRef<HTMLInputElement>(null);

  // Load recent searches from localStorage
  useEffect(() => {
    const saved = localStorage.getItem('recentSearches');
    if (saved) {
      setRecentSearches(JSON.parse(saved));
    }
  }, []);

  // Keyboard shortcut (Ctrl/Cmd + K)
  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === 'k' && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((open) => !open);
      }
    };

    document.addEventListener('keydown', down);
    return () => document.removeEventListener('keydown', down);
  }, []);

  const handleSearch = (searchQuery: string) => {
    if (searchQuery.trim()) {
      const trimmedQuery = searchQuery.trim();
      
      // Add to recent searches
      const updated = [trimmedQuery, ...recentSearches.filter(s => s !== trimmedQuery)].slice(0, 5);
      setRecentSearches(updated);
      localStorage.setItem('recentSearches', JSON.stringify(updated));
      
      navigate(`/search?q=${encodeURIComponent(trimmedQuery)}`);
      setOpen(false);
      setQuery('');
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    handleSearch(query);
  };

  const clearRecentSearches = () => {
    setRecentSearches([]);
    localStorage.removeItem('recentSearches');
  };

  return (
    <>
      {/* Desktop Search Bar */}
      <form onSubmit={handleSubmit} className="hidden sm:flex w-full max-w-sm items-center space-x-2">
        <div className="relative flex-1">
          <Input
            ref={inputRef}
            type="text"
            placeholder="Search novels... (⌘K)"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="input-field pr-10"
            onClick={() => setOpen(true)}
          />
          <div className="absolute right-3 top-1/2 transform -translate-y-1/2 text-xs text-muted-foreground border rounded px-1.5 py-0.5 bg-muted/50">
            ⌘K
          </div>
        </div>
        <Button type="submit" variant="outline" size="icon" className="shrink-0">
          <Search className="h-4 w-4" />
        </Button>
      </form>

      {/* Mobile Search Button */}
      <Button 
        variant="outline" 
        size="icon" 
        className="sm:hidden"
        onClick={() => setOpen(true)}
      >
        <Search className="h-4 w-4" />
      </Button>

      {/* Command Dialog */}
      <CommandDialog open={open} onOpenChange={setOpen}>
        <CommandInput 
          placeholder="Search novels..." 
          value={query}
          onValueChange={setQuery}
        />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          
          {recentSearches.length > 0 && (
            <CommandGroup heading={
              <div className="flex items-center justify-between">
                <span>Recent Searches</span>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={clearRecentSearches}
                  className="h-auto p-1 text-xs"
                >
                  <X className="h-3 w-3" />
                </Button>
              </div>
            }>
              {recentSearches.map((search, index) => (
                <CommandItem
                  key={index}
                  onSelect={() => handleSearch(search)}
                  className="cursor-pointer"
                >
                  <Clock className="mr-2 h-4 w-4 text-muted-foreground" />
                  {search}
                </CommandItem>
              ))}
            </CommandGroup>
          )}

          <CommandGroup heading="Quick Actions">
            <CommandItem onSelect={() => navigate('/add')}>
              <Search className="mr-2 h-4 w-4" />
              Add New Novel
            </CommandItem>
            <CommandItem onSelect={() => navigate('/library')}>
              <Search className="mr-2 h-4 w-4" />
              Browse Library
            </CommandItem>
          </CommandGroup>

          {query && (
            <CommandGroup heading="Search">
              <CommandItem onSelect={() => handleSearch(query)}>
                <Search className="mr-2 h-4 w-4" />
                Search for "{query}"
              </CommandItem>
            </CommandGroup>
          )}
        </CommandList>
      </CommandDialog>
    </>
  );
};

export default SearchBar;
