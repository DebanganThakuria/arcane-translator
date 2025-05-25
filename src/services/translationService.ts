// Standard library imports
import { Novel } from '../types/novel';
import { NovelDetails, TranslatedChapter } from '../types/api';

// API base URL - will be replaced by actual backend calls
const API_BASE_URL = 'http://localhost:8088';

// Extract and translate novel details from a webpage URL
export const extractNovelDetails = async (url: string, sourceSite: string): Promise<NovelDetails> => {
  try {
    // First, validate the URL format
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      throw new Error('Invalid URL format. URL must start with http:// or https://');
    }
    
    const response = await fetch(`${API_BASE_URL}/novels/extract`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ url: url, source: sourceSite })
    });
    
    if (!response.ok) {
      let errorMessage = 'Failed to fetch novel details';
      
      try {
        const errorData = await response.json();
        
        // Handle specific error codes with user-friendly messages
        if (response.status === 404) {
          errorMessage = 'Novel not found. Please check if the URL is correct.';
        } else if (response.status === 403) {
          errorMessage = 'Access to this website is forbidden. The website may be blocking our requests.';
        } else if (response.status === 500) {
          errorMessage = errorData.detail || 'Server error occurred while extracting novel details.';
        } else {
          errorMessage = errorData.detail || errorMessage;
        }
      } catch (e) {
        // If parsing JSON fails, use a generic message
        if (response.status === 404) {
          errorMessage = 'Novel not found. Please check if the URL is correct.';
        } else if (response.status === 403) {
          errorMessage = 'Access to this website is forbidden. The website may be blocking our requests.';
        }
      }
      
      throw new Error(errorMessage);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error extracting novel details:', error);
    throw error;
  }
};

export const addExtractedNovel = async (novelDetails: NovelDetails, sourceSite?: string, sourceUrl?: string): Promise<Novel> => {
  try {
    const novelPayload = {
      title: novelDetails.novel_title_translated,
      originalTitle: novelDetails.novel_title_original,
      cover: null, // Cover would need to be uploaded separately or handled differently
      source: sourceSite || "unknown",
      url: sourceUrl || "",
      summary: novelDetails.novel_summary_translated,
      author: novelDetails.novel_author_name_translated,
      status: novelDetails.status,
      genres: novelDetails.possible_novel_genres,
      chaptersCount: novelDetails.number_of_chapters
    };

    const response = await fetch(`${API_BASE_URL}/novels/add`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(novelPayload)
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ detail: 'Unknown error' }));
      throw new Error(errorData.detail || 'Failed to add novel');
    }

    return await response.json();
  } catch (error) {
    console.error('Error adding novel:', error);
    throw error;
  }
};

// Translate a chapter
export const translateChapter = async (
  novelId: string, 
  chapterNumber: number, 
  chapterUrl?: string, 
  originalContent?: string,
  saveToDb: boolean = false
): Promise<TranslatedChapter> => {
  try {
    // Build the query parameters for the POST request
    const url = new URL(`${API_BASE_URL}/chapters/translate`);
    url.searchParams.append("novel_id", novelId);
    url.searchParams.append("chapter_number", chapterNumber.toString());
    url.searchParams.append("save_to_db", saveToDb.toString());
    
    if (chapterUrl) url.searchParams.append("chapter_url", chapterUrl);
    if (originalContent) url.searchParams.append("original_content", originalContent);
    
    const response = await fetch(url.toString(), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ detail: 'Unknown error' }));
      throw new Error(errorData.detail || 'Failed to translate chapter');
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error translating chapter:', error);
    throw error;
  }
};

// Get all novels
export const getNovels = async (): Promise<Novel[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels`);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ detail: 'Unknown error' }));
      throw new Error(errorData.detail || 'Failed to fetch novels');
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error fetching novels:', error);
    throw error;
  }
};

// Get a single novel by ID
export const getNovel = async (novelId: string): Promise<Novel> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels/${novelId}`);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ detail: 'Unknown error' }));
      throw new Error(errorData.detail || 'Failed to fetch novel');
    }
    
    return await response.json();
  } catch (error) {
    console.error(`Error fetching novel with ID ${novelId}:`, error);
    throw error;
  }
};

// Refresh a novel (check for new chapters)
export const refreshNovel = async (novelId: string): Promise<{
  success: boolean;
  newChaptersCount: number;
  message: string;
  updatedDetails: {
    summary: string;
    author: string;
    status: string;
    genres: string[];
  };
}> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels/${novelId}/refresh`, {
      method: 'POST'
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ detail: 'Unknown error' }));
      throw new Error(errorData.detail || 'Failed to refresh novel');
    }
    
    return await response.json();
  } catch (error) {
    console.error(`Error refreshing novel with ID ${novelId}:`, error);
    throw error;
  }
};

// Set the first chapter URL for a novel
export const setFirstChapterUrl = async (
  novelId: string, 
  firstChapterUrl: string, 
  autoTranslate: boolean = true
): Promise<{
  success: boolean;
  message: string;
  firstChapterUrl?: string;
  urlPattern?: string;
  chapterTranslated?: boolean;
}> => {
  try {
    // Validate URL format
    if (!firstChapterUrl.startsWith('http://') && !firstChapterUrl.startsWith('https://')) {
      throw new Error('Invalid URL format. URL must start with http:// or https://');
    }

    const response = await fetch(
      `${API_BASE_URL}/novels/${novelId}/first-chapter?first_chapter_url=${encodeURIComponent(firstChapterUrl)}&auto_translate=${autoTranslate}`, 
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        }
      }
    );
    
    if (!response.ok) {
      let errorMessage = 'Failed to set first chapter URL';
      
      try {
        const errorData = await response.json();
        errorMessage = errorData.detail || errorMessage;
      } catch (e) {
        // If parsing JSON fails, use the generic message
      }
      
      throw new Error(errorMessage);
    }
    
    return await response.json();
  } catch (error) {
    console.error(`Error setting first chapter URL for novel ID ${novelId}:`, error);
    throw error;
  }
};

// Get the URL for a specific chapter
export const getChapterUrl = async (novelId: string, chapterNumber: number): Promise<string> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels/${novelId}/chapter-url/${chapterNumber}`);
    
    if (!response.ok) {
      let errorMessage = 'Failed to get chapter URL';
      
      try {
        const errorData = await response.json();
        errorMessage = errorData.detail || errorMessage;
      } catch (e) {
        // If parsing JSON fails, use the generic message
      }
      
      throw new Error(errorMessage);
    }
    
    const result = await response.json();
    return result.url;
  } catch (error) {
    console.error(`Error getting chapter URL for novel ID ${novelId}, chapter ${chapterNumber}:`, error);
    throw error;
  }
};

export default {
  extractNovelDetails,
  addExtractedNovel,
  translateChapter,
  getNovels,
  getNovel,
  refreshNovel,
  setFirstChapterUrl,
  getChapterUrl
};
