// Standard library imports
import {Chapter, Novel} from '../types/novel';

// API base URL - will be replaced by actual backend calls
const API_BASE_URL = 'http://localhost:8088';

// Extract and translate novel details from a webpage URL
export const extractNovelDetails = async (url: string, sourceSite: string): Promise<Novel> => {
  try {
    // First, validate the URL format
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      throw new Error('Invalid URL format. URL must start with http:// or https://');
    }
    
    const response = await fetch(`${API_BASE_URL}/novels/translate`, {
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

// Translate a chapter
export const translateChapter = async (
  novelId: string,
): Promise<Chapter> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels/translate/chapter`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ novel_id: novelId})
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

// Set the first chapter URL for a novel
export const setFirstChapterUrl = async (
  novelId: string, 
  firstChapterUrl: string, 
  autoTranslate: boolean = true
): Promise<Chapter> => {
  try {
    // Validate URL format
    if (!firstChapterUrl.startsWith('http://') && !firstChapterUrl.startsWith('https://')) {
      throw new Error('Invalid URL format. URL must start with http:// or https://');
    }

    const response = await fetch(
      `${API_BASE_URL}/novels/translate/first_chapter`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ novel_id: novelId, chapter_url: firstChapterUrl })
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

export default {
  extractNovelDetails,
  translateChapter,
  setFirstChapterUrl,
};
