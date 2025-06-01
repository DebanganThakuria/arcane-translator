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
  last_read_chapter_number: number;
  last_read_timestamp: number;
  last_updated: number;
  date_added: number;
}

export interface PaginatedNovels {
  novels: Novel[];
  total_count: number;
  current_page: number;
  total_pages: number;
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
