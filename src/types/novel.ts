
export interface Novel {
  id: string;
  title: string;
  originalTitle?: string;
  cover?: string;
  source: string;
  url: string;
  summary: string;
  author?: string;
  status?: 'Ongoing' | 'Completed' | 'Unknown';
  genres?: string[];
  chaptersCount: number;
  firstChapterUrl?: string;
  urlPattern?: string;  // Pattern for generating chapter URLs
  lastRead?: {
    chapterId: string;
    progress: number;
    timestamp: number;
  };
  lastUpdated: number;
  dateAdded: number;
}

export interface Chapter {
  id: string;
  novelId: string;
  number: number;
  title: string;
  originalTitle?: string;
  content: string;
  dateTranslated: number;
  wordCount?: number;
  url?: string;              // URL of this chapter
  nextChapterUrl?: string;   // URL to the next chapter
}

export interface SourceSite {
  id: string;
  name: string;
  url: string;
  language: 'Chinese' | 'Korean' | 'Japanese' | 'Other';
  icon?: string;
}
