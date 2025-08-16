import { Novel, PaginatedNovels, Chapter, SourceSite } from '../types/novel';
import {NovelStats} from "@/types/api.ts";

// API base URL - will be replaced by actual backend calls
const API_BASE_URL = 'http://localhost:8088';

// Fetch novels stats
export const getNovelsStats = async () => {
    try {
        const response = await fetch(`${API_BASE_URL}/stats/novels`);
        if (!response.ok) throw new Error('Failed to fetch novels stats');
        const data = await response.json();
        return {
            novel_count: data.novel_count || 0,
            chapter_count: data.chapter_count || 0,
            languages_supported: 3 // This is hardcoded as it's not coming from the backend
        };
    } catch (error) {
        console.error('Error fetching novels stats:', error);
        return {
            novel_count: 0,
            chapter_count: 0,
            languages_supported: 3
        };
    }
};

// Fetch all novels with pagination
export const getAllNovels = async (page = 1, limit = 20): Promise<PaginatedNovels> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels?page=${page}&limit=${limit}`);
    if (!response.ok) throw new Error('Failed to fetch novels');
    return await response.json();
  } catch (error) {
    console.error('Error fetching novels:', error);
    return {
      novels: [],
      total_count: 0,
      current_page: page,
      total_pages: 1
    };
  }
};

// Fetch novels by filter from backend
export const getNovelsByFilter = async (filter: string, value: string, page = 1, limit = 20): Promise<PaginatedNovels> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels?filter=${filter}&value=${value}&page=${page}&limit=${limit}`);
    if (!response.ok) throw new Error('Failed to fetch novels');
    return await response.json();
  } catch (error) {
    console.error('Error fetching novels:', error);
    return {
      novels: [],
      total_count: 0,
      current_page: page,
      total_pages: 1
    };
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

// Get chapter by number
export const getChapterByNumber = async (novelId: string, chapterNumber: number): Promise<Chapter | undefined> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels/${novelId}/chapters/num/${chapterNumber}`);
    if (!response.ok) {
      console.error(`Failed to fetch chapter ${chapterNumber} for novel ${novelId}`);
      return undefined;
    }
    return await response.json();
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

// Search novels
export const searchNovelsByQuery = async (query: string): Promise<Novel[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/search/novels/${query}`);
    if (!response.ok) throw new Error('Failed to search novels');
    return await response.json();
  } catch (error) {
    console.error('Error searching novels:', error);
    return [];
  }
};

// Delete novel
export const deleteNovel = async (id: string): Promise<void> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete novel');
  } catch (error) {
    console.error('Error deleting novel:', error);
    throw error;
  }
};

// Delete chapter
export const deleteChapter = async (novelId: string, chapterId: string): Promise<void> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels/${novelId}/chapters/${chapterId}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete chapter');
  } catch (error) {
    console.error('Error deleting chapter:', error);
    throw error;
  }
};
