
import { getNovelById } from '../database/db';
import { refreshNovel } from './translationService';

// Function to refresh novel data using the backend API
export const refreshNovelData = async (novelId: string): Promise<{ 
  success: boolean;
  message: string; 
  newChaptersCount?: number;
  updatedDetails?: {
    summary: string;
    author: string;
    status: string;
    genres: string[];
  };
}> => {
  try {
    // Get current novel data to make sure it exists
    const novel = await getNovelById(novelId);
    
    if (!novel) {
      return {
        success: false,
        message: "Novel not found"
      };
    }
    
    // Call the refreshNovel function from the translation service
    const refreshResult = await refreshNovel(novelId);
    
    console.log(`Refreshed novel ${novel.title}, found ${refreshResult.newChaptersCount || 0} new chapters`);
    
    return {
      success: refreshResult.success,
      message: refreshResult.message,
      newChaptersCount: refreshResult.newChaptersCount || 0,
      updatedDetails: refreshResult.updatedDetails
    };
  } catch (error) {
    console.error("Error refreshing novel data:", error);
    const errorMessage = error instanceof Error ? error.message : "Failed to refresh novel data";
    return {
      success: false,
      message: errorMessage
    };
  }
};
