
import { Novel, Chapter, SourceSite } from '../types/novel';
import { mockNovels, mockChapters, mockSites } from '../data/mockData';

// Simple in-memory database functions that will be replaced by Python backend API calls
// These are temporary fallbacks until the Python backend is implemented

export const getAllNovels = (): Novel[] => {
  return mockNovels;
};

export const getNovelById = (id: string): Novel | undefined => {
  return mockNovels.find(novel => novel.id === id);
};

export const getChaptersForNovel = (novelId: string): Chapter[] => {
  return mockChapters.filter(chapter => chapter.novelId === novelId);
};

export const getChapter = (novelId: string, chapterNumber: number): Chapter | undefined => {
  return mockChapters.find(
    chapter => chapter.novelId === novelId && chapter.number === chapterNumber
  );
};

export const getAllSourceSites = (): SourceSite[] => {
  return mockSites;
};

// These functions will be replaced by Python backend API calls
export const addNovel = async (novel: Novel): Promise<boolean> => {
  console.log('Adding novel (will be handled by Python backend):', novel);
  return true;
};

export const updateNovel = async (id: string, updates: Partial<Novel>): Promise<boolean> => {
  console.log('Updating novel (will be handled by Python backend):', id, updates);
  return true;
};

export const addChapter = async (chapter: Chapter): Promise<boolean> => {
  console.log('Adding chapter (will be handled by Python backend):', chapter);
  return true;
};
