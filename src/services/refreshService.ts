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
    const refreshResult = await refreshNovel(novelId);
    
    console.log(`Refreshed novel and found ${refreshResult.newChaptersCount || 0} new chapters`);
    
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
