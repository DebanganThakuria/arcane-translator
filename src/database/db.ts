
import { Novel, Chapter, SourceSite } from '../types/novel';

// API base URL - will be replaced by actual backend calls
const API_BASE_URL = 'http://localhost:8088';

// Fetch all novels from backend
export const getAllNovels = async (): Promise<Novel[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels`);
    if (!response.ok) throw new Error('Failed to fetch novels');
    return await response.json();
  } catch (error) {
    console.error('Error fetching novels:', error);
    return [];
  }
};

// Get novel by ID from backend
export const getNovelById = async (id: string): Promise<Novel | undefined> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels/${id}`);
    if (!response.ok) return undefined;
    return await response.json();
  } catch (error) {
    console.error('Error fetching novel:', error);
    return undefined;
  }
};

// Get chapters for a novel from backend
export const getChaptersForNovel = async (novelId: string): Promise<Chapter[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels/${novelId}/chapters`);
    if (!response.ok) throw new Error('Failed to fetch chapters');
    return await response.json();
  } catch (error) {
    console.error('Error fetching chapters:', error);
    return [];
  }
};

// Get specific chapter from backend
export const getChapter = async (novelId: string, chapterNumber: number): Promise<Chapter | undefined> => {
  try {
    const chapters = await getChaptersForNovel(novelId);
    return chapters.find(chapter => chapter.number === chapterNumber);
  } catch (error) {
    console.error('Error fetching chapter:', error);
    return undefined;
  }
};

// Get all source sites from backend
export const getAllSourceSites = async (): Promise<SourceSite[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/sources`);
    if (!response.ok) throw new Error('Failed to fetch sources');
    return await response.json();
  } catch (error) {
    console.error('Error fetching sources:', error);
    return [];
  }
};

// Add novel via backend
export const addNovel = async (novel: Novel): Promise<boolean> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(novel)
    });
    return response.ok;
  } catch (error) {
    console.error('Error adding novel:', error);
    return false;
  }
};

// Update novel via backend
export const updateNovel = async (id: string, updates: Partial<Novel>): Promise<boolean> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(updates)
    });
    return response.ok;
  } catch (error) {
    console.error('Error updating novel:', error);
    return false;
  }
};

// Add chapter via backend
export const addChapter = async (chapter: Chapter): Promise<boolean> => {
  try {
    const response = await fetch(`${API_BASE_URL}/chapters`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(chapter)
    });
    return response.ok;
  } catch (error) {
    console.error('Error adding chapter:', error);
    return false;
  }
};

// Save or update chapter via backend
export const saveChapter = async (chapter: Chapter): Promise<boolean> => {
  try {
    // Check if chapter already exists
    const existingChapter = await getChapter(chapter.novel_id, chapter.number);
    
    if (existingChapter) {
      // Update existing chapter
      const response = await fetch(`${API_BASE_URL}/chapters/${chapter.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(chapter)
      });
      return response.ok;
    } else {
      // Add new chapter
      return await addChapter(chapter);
    }
  } catch (error) {
    console.error('Error saving chapter:', error);
    return false;
  }
};
