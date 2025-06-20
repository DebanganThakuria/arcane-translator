import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { ThemeProvider } from "next-themes";
import Index from "./pages/Index";
import Library from "./pages/Library";
import NovelDetail from "./pages/NovelDetail";
import ChapterReader from "./pages/ChapterReader";
import AddNovel from "./pages/AddNovel";
import SearchResults from "./pages/SearchResults";
import NotFound from "./pages/NotFound";
import ChineseNovels from "./pages/ChineseNovels";
import KoreanNovels from "./pages/KoreanNovels";
import JapaneseNovels from "./pages/JapaneseNovels";
import GenrePage from "./pages/GenrePage";
import { sourceManager } from "./utils/sourceManager";
import { useEffect } from "react";

const queryClient = new QueryClient();

const App = () => {
  // Preload source sites when app starts
  useEffect(() => {
    sourceManager.loadSources().catch(error => {
      console.error('Failed to preload source sites:', error);
    });
  }, []);

  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/novel/:novelId/chapter/:chapterNumber" element={
            <ThemeProvider attribute="class" forcedTheme="light" disableTransitionOnChange>
              <TooltipProvider>
                <Toaster />
                <Sonner />
                <ChapterReader />
              </TooltipProvider>
            </ThemeProvider>
          } />
          <Route path="*" element={
            <ThemeProvider attribute="class" defaultTheme="light" enableSystem disableTransitionOnChange>
              <TooltipProvider>
                <Toaster />
                <Sonner />
                <Routes>
                  <Route path="/" element={<Index />} />
                  <Route path="/library" element={<Library />} />
                  <Route path="/novel/:id" element={<NovelDetail />} />
                  <Route path="/add" element={<AddNovel />} />
                  <Route path="/search" element={<SearchResults />} />
                  <Route path="/chinese" element={<ChineseNovels />} />
                  <Route path="/korean" element={<KoreanNovels />} />
                  <Route path="/japanese" element={<JapaneseNovels />} />
                  <Route path="/genre/:genre" element={<GenrePage />} />
                  <Route path="*" element={<NotFound />} />
                </Routes>
              </TooltipProvider>
            </ThemeProvider>
          } />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
};

export default App;
