interface ReadingProgressData {
  lastChapter: number;
  lastReadAt: number;
  progress: number; // percentage through the chapter (0-100)
  chapterTitle?: string;
}

interface ReadingProgressStorage {
  [novelId: string]: ReadingProgressData;
}

const STORAGE_KEY = 'novel_reading_progress';

export const getReadingProgress = (novelId: string): ReadingProgressData | null => {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (!stored) return null;
    
    const progressData: ReadingProgressStorage = JSON.parse(stored);
    return progressData[novelId] || null;
  } catch (error) {
    console.error('Error reading progress from localStorage:', error);
    return null;
  }
};

export const saveReadingProgress = (
  novelId: string, 
  chapterNumber: number, 
  progress: number = 0,
  chapterTitle?: string
): void => {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    const progressData: ReadingProgressStorage = stored ? JSON.parse(stored) : {};
    
    progressData[novelId] = {
      lastChapter: chapterNumber,
      lastReadAt: Date.now(),
      progress: Math.min(Math.max(progress, 0), 100), // Clamp between 0-100
      chapterTitle
    };
    
    localStorage.setItem(STORAGE_KEY, JSON.stringify(progressData));
  } catch (error) {
    console.error('Error saving progress to localStorage:', error);
  }
};

export const removeReadingProgress = (novelId: string): void => {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (!stored) return;
    
    const progressData: ReadingProgressStorage = JSON.parse(stored);
    delete progressData[novelId];
    
    localStorage.setItem(STORAGE_KEY, JSON.stringify(progressData));
  } catch (error) {
    console.error('Error removing progress from localStorage:', error);
  }
};

export const getAllReadingProgress = (): ReadingProgressStorage => {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    return stored ? JSON.parse(stored) : {};
  } catch (error) {
    console.error('Error reading all progress from localStorage:', error);
    return {};
  }
};

export const hasReadingProgress = (novelId: string): boolean => {
  const progress = getReadingProgress(novelId);
  return progress !== null && progress.lastChapter > 0;
};

export const getReadingProgressSummary = (novelId: string): string | null => {
  const progress = getReadingProgress(novelId);
  if (!progress) return null;
  
  const timeAgo = new Date(progress.lastReadAt);
  const now = new Date();
  const diffInHours = Math.floor((now.getTime() - timeAgo.getTime()) / (1000 * 60 * 60));
  const diffInMinutes = Math.floor((now.getTime() - timeAgo.getTime()) / (1000 * 60));
  
  let timeStr = '';
  if (diffInMinutes < 1) {
    timeStr = 'just now';
  } else if (diffInMinutes < 60) {
    timeStr = `${diffInMinutes} minute${diffInMinutes > 1 ? 's' : ''} ago`;
  } else if (diffInHours < 24) {
    timeStr = `${diffInHours} hour${diffInHours > 1 ? 's' : ''} ago`;
  } else {
    const diffInDays = Math.floor(diffInHours / 24);
    timeStr = `${diffInDays} day${diffInDays > 1 ? 's' : ''} ago`;
  }
  
  return `Chapter ${progress.lastChapter} • ${Math.round(progress.progress)}% • ${timeStr}`;
}; 