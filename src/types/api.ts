export interface NovelDetails {
  novel_title_original: string
  novel_title_translated: string;
  novel_summary_translated: string;
  novel_author_name_translated?: string;
  possible_novel_genres?: string[];
  number_of_chapters: number;
  status?: 'Ongoing' | 'Completed' | 'Unknown';
}

export interface TranslatedChapter {
  translated_chapter_title: string;
  translated_chapter_contents: string;
  possible_new_genres?: string[];
  word_count?: number;
}