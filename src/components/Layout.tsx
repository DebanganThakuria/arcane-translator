
import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import SearchBar from './SearchBar';
import { Book, Library, Plus } from 'lucide-react';

interface LayoutProps {
  children: React.ReactNode;
  hideNavigation?: boolean;
}

const Layout: React.FC<LayoutProps> = ({ children, hideNavigation = false }) => {
  const location = useLocation();
  const pathname = location.pathname;
  
  return (
    <div className="min-h-screen flex flex-col">
      {!hideNavigation && (
        <header className="border-b">
          <div className="container mx-auto py-4 px-4 flex flex-col sm:flex-row items-center justify-between space-y-4 sm:space-y-0">
            <Link to="/" className="text-2xl font-bold text-novel">
              Arcane Translator
            </Link>
            
            <div className="flex items-center space-x-4">
              <SearchBar />
              <Link to="/add">
                <Button size="sm" variant="outline">
                  <Plus className="h-4 w-4 mr-2" />
                  Add Novel
                </Button>
              </Link>
            </div>
          </div>
          
          <div className="container mx-auto px-4 pb-2">
            <Tabs value={pathname === "/" ? "/" : pathname.split("/")[1]} className="w-full">
              <TabsList className="w-full justify-start">
                <TabsTrigger value="/" asChild>
                  <Link to="/">
                    <Book className="h-4 w-4 mr-2" />
                    Home
                  </Link>
                </TabsTrigger>
                <TabsTrigger value="library" asChild>
                  <Link to="/library">
                    <Library className="h-4 w-4 mr-2" />
                    Library
                  </Link>
                </TabsTrigger>
              </TabsList>
            </Tabs>
          </div>
        </header>
      )}
      
      <main className="flex-1">
        {children}
      </main>
      
      {!hideNavigation && (
        <footer className="border-t py-6">
          <div className="container mx-auto px-4 text-center text-sm text-muted-foreground">
            &copy; {new Date().getFullYear()} Arcane Translator • AI-Powered Webnovel Translation
          </div>
        </footer>
      )}
    </div>
  );
};

export default Layout;
