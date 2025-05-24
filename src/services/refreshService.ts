
import { getNovelById } from '../database/db';

// Function to refresh novel data - calls Python backend API
export const refreshNovelData = async (novelId: string): Promise<{ 
  success: boolean;
  message: string; 
  newChaptersCount?: number;
}> => {
  try {
    // Get current novel data
    const novel = await getNovelById(novelId);
    
    if (!novel) {
      return {
        success: false,
        message: "Novel not found"
      };
    }
    
    // Call Python backend API to refresh novel data
    const response = await fetch(`http://localhost:8000/novels/${novelId}/refresh`, {
      method: 'POST'
    });
    
    if (!response.ok) {
      throw new Error('Failed to refresh novel data');
    }
    
    const result = await response.json();
    
    console.log(`Refreshed novel ${novel.title}, found ${result.newChaptersCount || 0} new chapters`);
    
    return {
      success: true,
      message: result.newChaptersCount > 0 
        ? `Found ${result.newChaptersCount} new chapter${result.newChaptersCount > 1 ? 's' : ''}!`
        : "No new chapters found",
      newChaptersCount: result.newChaptersCount || 0
    };
  } catch (error) {
    console.error("Error refreshing novel data:", error);
    return {
      success: false,
      message: "Failed to refresh novel data"
    };
  }
};
