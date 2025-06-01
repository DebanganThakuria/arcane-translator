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
        <header className="border-b bg-white/95 backdrop-blur-md shadow-sm">
          <div className="container mx-auto py-4 px-4 flex flex-col sm:flex-row items-center justify-between space-y-4 sm:space-y-0">
            <Link to="/" className="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">
              Arcane Translator
            </Link>
            
            <div className="flex items-center space-x-4">
              <SearchBar />
              <Link to="/add">
                <Button size="sm" variant="outline" className="border-indigo-200 hover:bg-indigo-50 text-indigo-600">
                  <Plus className="h-4 w-4 mr-2" />
                  Add Novel
                </Button>
              </Link>
            </div>
          </div>
          
          <div className="container mx-auto px-4 pb-2">
            <Tabs value={pathname === "/" ? "/" : pathname.split("/")[1]} className="w-full">
              <TabsList className="w-full justify-start bg-white/70 backdrop-blur-sm">
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
                <TabsTrigger value="chinese" asChild>
                  <Link to="/chinese">
                    Chinese
                  </Link>
                </TabsTrigger>
                <TabsTrigger value="korean" asChild>
                  <Link to="/korean">
                    Korean
                  </Link>
                </TabsTrigger>
                <TabsTrigger value="japanese" asChild>
                  <Link to="/japanese">
                    Japanese
                  </Link>
                </TabsTrigger>
              </TabsList>
            </Tabs>
          </div>
        </header>
      )}
      
      <main className="flex-1 app-container">
        {children}
      </main>
      
      {!hideNavigation && (
        <footer className="border-t py-6 bg-white/80 backdrop-blur-sm">
          <div className="container mx-auto px-4 text-center text-sm text-muted-foreground">
            &copy; {new Date().getFullYear()} Arcane Translator • AI-Powered Webnovel Translation
          </div>
        </footer>
      )}
    </div>
  );
};

export default Layout;
