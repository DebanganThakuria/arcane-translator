import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useToast } from '@/components/ui/use-toast';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogClose,
} from '@/components/ui/dialog';
import { setFirstChapterUrl } from '../services/translationService';

interface FirstChapterDialogProps {
  novelId: string;
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

const FirstChapterDialog = ({ novelId, isOpen, onOpenChange, onSuccess }: FirstChapterDialogProps) => {
  const [url, setUrl] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { toast } = useToast();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!url.trim()) {
      toast({
        title: "Error",
        description: "Please enter a valid URL",
        variant: "destructive",
      });
      return;
    }
    
    setIsSubmitting(true);
    
    try {
      const toastId = toast({
        title: "Setting first chapter URL",
      }).id;
      
      const result = await setFirstChapterUrl(novelId, url);
      
      toast({
        title: "Success",
        description: result.chapterTranslated 
          ? "First chapter URL set and chapter translated successfully. You can now start reading." 
          : "First chapter URL set successfully. You can now start reading.",
      });
      
      setUrl('');
      onOpenChange(false);
      onSuccess();
    } catch (error) {
      console.error('Error setting first chapter URL:', error);
      const errorMessage = error instanceof Error ? error.message : 'Unknown error';
      toast({
        title: "Error",
        description: `Failed to set first chapter URL: ${errorMessage}`,
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Set First Chapter URL</DialogTitle>
          <DialogDescription>
            Enter the URL of the first chapter to enable sequential reading.
            <br />
            For example: https://www.69shuba.com/txt/84522/39367190
          </DialogDescription>
        </DialogHeader>
        
        <form onSubmit={handleSubmit} className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="chapterUrl">Chapter URL</Label>
            <Input
              id="chapterUrl"
              placeholder="https://example.com/novel/12345/chapter/1"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              required
            />
            <p className="text-xs text-muted-foreground">
              Enter the complete URL to the first chapter of the novel.
              Next chapters will be automatically detected using the pattern.
            </p>
          </div>
          
          <DialogFooter className="sm:justify-between">
            <DialogClose asChild>
              <Button type="button" variant="outline" disabled={isSubmitting}>
                Cancel
              </Button>
            </DialogClose>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? 'Setting...' : 'Set First Chapter'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default FirstChapterDialog;
