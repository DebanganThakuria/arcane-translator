import { Novel, PaginatedNovels, Chapter, SourceSite } from '../types/novel';

// API base URL - will be replaced by actual backend calls
const API_BASE_URL = 'http://localhost:8088';

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
    const response = await fetch(`${API_BASE_URL}/novels/search/${query}`);
    if (!response.ok) throw new Error('Failed to search novels');
    return await response.json();
  } catch (error) {
    console.error('Error searching novels:', error);
    return [];
  }
};
