
import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { useToast } from '@/components/ui/use-toast';
import { getAllSourceSites } from '../database/db';
import { SourceSite } from '../types/novel';
import { extractNovelDetails, isGeminiConfigured, setGeminiApiKey } from '../services/translationService';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';

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

  useEffect(() => {
    // Check if API key is configured, if not, show dialog
    if (!isGeminiConfigured()) {
      setIsApiKeyDialogOpen(true);
    }
  }, []);

  const handleApiKeySubmit = () => {
    if (!apiKey.trim()) {
      toast({
        title: "Error",
        description: "Please enter a valid API key",
        variant: "destructive",
      });
      return;
    }

    setGeminiApiKey(apiKey);
    setIsApiKeyDialogOpen(false);
    toast({
      title: "Success",
      description: "API key configured successfully",
    });
  };

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

    if (!isGeminiConfigured()) {
      setIsApiKeyDialogOpen(true);
      return;
    }
    
    setIsSubmitting(true);
    
    try {
      toast({
        title: "Processing",
        description: "Extracting and translating novel information...",
      });

      const novelDetails = await extractNovelDetails(url, source);
      
      toast({
        title: "Success",
        description: `Novel "${novelDetails.novel_title_translated || 'Unknown Title'}" added successfully!`,
      });
      
      setTimeout(() => {
        setIsSubmitting(false);
        navigate('/library');
      }, 2000);
    } catch (error) {
      console.error('Error processing novel:', error);
      toast({
        title: "Error",
        description: "Failed to process novel. Please check the URL and try again.",
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
          <DialogFooter>
            <Button type="button" onClick={handleApiKeySubmit}>Save API Key</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
};

export default AddNovelForm;
