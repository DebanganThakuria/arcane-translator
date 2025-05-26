
import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { useToast } from '@/components/ui/use-toast';
import { getAllSourceSites } from '../database/db';
import { SourceSite } from '../types/novel';
import { extractNovelDetails } from '../services/translationService';
import { NovelDetails } from '../types/api';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';

const AddNovelForm = () => {
  const [url, setUrl] = useState('');
  const [source, setSource] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [sites, setSites] = useState<SourceSite[]>([]);
  const [isApiKeyDialogOpen, setIsApiKeyDialogOpen] = useState(false);
  const [apiKey, setApiKey] = useState('');
  const [isLoadingSites, setIsLoadingSites] = useState(true);
  const { toast } = useToast();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchSites = async () => {
      try {
        setIsLoadingSites(true);
        const fetchedSites = await getAllSourceSites();
        setSites(fetchedSites || []);
      } catch (error) {
        console.error('Error fetching source sites:', error);
        setSites([]);
        toast({
          title: "Error",
          description: "Failed to load source sites. Using fallback options.",
          variant: "destructive",
        });
      } finally {
        setIsLoadingSites(false);
      }
    };

    fetchSites();
  }, [toast]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!url.trim() || !source.trim()) {
      toast({
        title: "Error",
        description: "Please fill in all fields",
        variant: "destructive",
      });
      return;
    }
    
    setIsSubmitting(true);
    
    try {
      toast({
        title: "Processing",
        description: "Extracting novel information...",
      });

      const novel = await extractNovelDetails(url, source);
      
      toast({
        title: "Extracting",
        description: `Found novel: "${novel.title}". Adding to library...`,
      });
      
      toast({
        title: "Success",
        description: `Novel "${novel.title}" added successfully!`,
      });
      
      // Navigate to the library after successful addition
      setTimeout(() => {
        setIsSubmitting(false);
        navigate('/library');
      }, 2000);
    } catch (error) {
      console.error('Error processing novel:', error);
      const errorMessage = error instanceof Error ? error.message : 'Unknown error';
      
      // Handle specific error messages with more user-friendly explanations
      let title = "Error";
      let description = `Failed to process novel: ${errorMessage}`;
      
      if (errorMessage.includes('404') || errorMessage.includes('not found')) {
        title = "Novel Not Found";
        description = "The URL you provided doesn't seem to contain a valid novel. Please check the URL and try again. Try using the novel's main page URL instead of a chapter page.";
      } else if (errorMessage.includes('403') || errorMessage.includes('forbidden')) {
        title = "Access Denied";
        description = "The website is blocking our access. This is common with some novel sites. Try another source or URL, or try a different novel site. Some sites actively block web scrapers.";
      } else if (errorMessage.includes('connection') || errorMessage.includes('network')) {
        title = "Connection Issue";
        description = "There was a problem connecting to the novel site. Please check your internet connection and try again. If the problem persists, the site might be down or blocking access.";
      }
      
      toast({
        title: title,
        description: description,
        variant: "destructive",
      });
      setIsSubmitting(false);
    }
  };

  const fallbackSites = [
    { id: 'qidian', name: 'Qidian', url: 'https://www.qidian.com', language: 'Chinese' as const },
    { id: 'naver', name: 'Naver Series', url: 'https://series.naver.com', language: 'Korean' as const },
    { id: 'syosetu', name: 'Syosetu', url: 'https://syosetu.com', language: 'Japanese' as const }
  ];

  const sitesToShow = sites.length > 0 ? sites : fallbackSites;

  return (
    <>
      <form onSubmit={handleSubmit} className="space-y-6 max-w-md mx-auto">
        <div className="space-y-2">
          <Label htmlFor="url">Novel URL</Label>
          <Input
            id="url"
            placeholder="https://example.com/novel/12345"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            required
          />
          <p className="text-xs text-muted-foreground">
            Enter the URL of the novel's main page
          </p>
        </div>
        
        <div className="space-y-2">
          <Label htmlFor="source">Source Site</Label>
          <Select value={source} onValueChange={setSource} required disabled={isLoadingSites}>
            <SelectTrigger>
              <SelectValue placeholder={isLoadingSites ? "Loading sources..." : "Select source site"} />
            </SelectTrigger>
            <SelectContent>
              {sitesToShow.map(site => (
                <SelectItem key={site.id} value={site.id}>
                  {site.name} ({site.language})
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        
        <Button type="submit" className="w-full" disabled={isSubmitting || isLoadingSites}>
          {isSubmitting ? "Processing..." : "Add Novel"}
        </Button>
      </form>

      <Dialog open={isApiKeyDialogOpen} onOpenChange={setIsApiKeyDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Google Generative AI API Key</DialogTitle>
            <DialogDescription>
              Please enter your Gemini API key to enable novel translations
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="apiKey">API Key</Label>
              <Input
                id="apiKey"
                placeholder="AIzaSyA1234BcDeFgHijKLmn..."
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
              />
              <p className="text-xs text-muted-foreground">
                Get your API key from{" "}
                <a 
                  href="https://makersuite.google.com/app/apikey" 
                  target="_blank" 
                  rel="noreferrer"
                  className="text-indigo-600 hover:text-indigo-800"
                >
                  Google AI Studio
                </a>
              </p>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
};

export default AddNovelForm;
