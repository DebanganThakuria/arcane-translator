
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { useToast } from '@/components/ui/use-toast';
import { mockSites } from '../data/mockData';

const AddNovelForm = () => {
  const [url, setUrl] = useState('');
  const [source, setSource] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { toast } = useToast();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!url || !source) {
      toast({
        title: "Error",
        description: "Please fill in all fields",
        variant: "destructive",
      });
      return;
    }
    
    setIsSubmitting(true);
    
    // Simulate API call to add a novel
    try {
      // In a real app, this would be an API call to add the novel
      await new Promise(resolve => setTimeout(resolve, 1500));
      
      toast({
        title: "Success",
        description: "Novel added successfully! Translation in progress...",
      });
      
      navigate('/library');
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to add novel. Please try again.",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
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
        <Select value={source} onValueChange={setSource} required>
          <SelectTrigger>
            <SelectValue placeholder="Select source site" />
          </SelectTrigger>
          <SelectContent>
            {mockSites.map(site => (
              <SelectItem key={site.id} value={site.id}>
                {site.name} ({site.language})
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      
      <Button type="submit" className="w-full" disabled={isSubmitting}>
        {isSubmitting ? "Adding Novel..." : "Add Novel"}
      </Button>
    </form>
  );
};

export default AddNovelForm;
