
export interface Novel {
  id: string;
  title: string;
  original_title?: string;
  cover?: string;
  source: string;
  url: string;
  summary: string;
  author?: string;
  status?: 'Ongoing' | 'Completed' | 'Unknown';
  genres?: string[];
  chapters_count: number;
  first_chapter_url?: string;
  url_pattern?: string;  // Pattern for generating chapter URLs
  last_updated: number;
  date_added: number;
}

export interface Chapter {
  id: string;
  novel_id: string;
  number: number;
  title: string;
  original_tile?: string;
  content: string;
  date_translated: number;
  word_count?: number;
  url?: string;              // URL of this chapter
}

export interface SourceSite {
  id: string;
  name: string;
  url: string;
  language: 'Chinese' | 'Korean' | 'Japanese' | 'Other';
  icon?: string;
}
