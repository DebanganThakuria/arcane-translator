
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
}

export interface SourceSite {
  id: string;
  name: string;
  url: string;
  language: 'Chinese' | 'Korean' | 'Japanese' | 'Other';
  icon?: string;
}
