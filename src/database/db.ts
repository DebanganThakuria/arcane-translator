
import Database from 'better-sqlite3';
import fs from 'fs';
import path from 'path';
import { Novel, Chapter, SourceSite } from '../types/novel';
import { mockNovels, mockChapters, mockSites } from '../data/mockData';

// Database file path
const DB_PATH = path.join(process.cwd(), 'novels.db');

// Initialize database connection
export const getDb = () => {
  const db = new Database(DB_PATH);
  db.pragma('journal_mode = WAL');
  return db;
};

// Initialize the database schema and seed with mock data
export const initializeDatabase = () => {
  // Check if database already exists
  const dbExists = fs.existsSync(DB_PATH);
  
  const db = getDb();
  
  // Create tables if they don't exist
  db.exec(`
    CREATE TABLE IF NOT EXISTS sources (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      url TEXT NOT NULL,
      language TEXT NOT NULL,
      icon TEXT
    );
    
    CREATE TABLE IF NOT EXISTS novels (
      id TEXT PRIMARY KEY,
      title TEXT NOT NULL,
      originalTitle TEXT,
      cover TEXT,
      source TEXT NOT NULL,
      url TEXT NOT NULL,
      summary TEXT NOT NULL,
      author TEXT,
      status TEXT,
      chaptersCount INTEGER NOT NULL DEFAULT 0,
      lastRead_chapterId TEXT,
      lastRead_progress REAL,
      lastRead_timestamp INTEGER,
      lastUpdated INTEGER NOT NULL,
      dateAdded INTEGER NOT NULL,
      FOREIGN KEY (source) REFERENCES sources(id)
    );
    
    CREATE TABLE IF NOT EXISTS novel_genres (
      novelId TEXT NOT NULL,
      genre TEXT NOT NULL,
      PRIMARY KEY (novelId, genre),
      FOREIGN KEY (novelId) REFERENCES novels(id)
    );
    
    CREATE TABLE IF NOT EXISTS chapters (
      id TEXT PRIMARY KEY,
      novelId TEXT NOT NULL,
      number INTEGER NOT NULL,
      title TEXT NOT NULL,
      originalTitle TEXT,
      content TEXT NOT NULL,
      dateTranslated INTEGER NOT NULL,
      wordCount INTEGER,
      FOREIGN KEY (novelId) REFERENCES novels(id)
    );
  `);
  
  // Seed database with mock data if it's a new database
  if (!dbExists) {
    // Insert sources
    const insertSource = db.prepare(`
      INSERT INTO sources (id, name, url, language, icon)
      VALUES (?, ?, ?, ?, ?)
    `);
    
    mockSites.forEach(site => {
      insertSource.run(
        site.id,
        site.name,
        site.url,
        site.language,
        site.icon
      );
    });
    
    // Insert novels
    const insertNovel = db.prepare(`
      INSERT INTO novels (
        id, title, originalTitle, cover, source, url, summary, author, 
        status, chaptersCount, lastRead_chapterId, lastRead_progress, 
        lastRead_timestamp, lastUpdated, dateAdded
      )
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `);
    
    const insertGenre = db.prepare(`
      INSERT INTO novel_genres (novelId, genre)
      VALUES (?, ?)
    `);
    
    mockNovels.forEach(novel => {
      insertNovel.run(
        novel.id,
        novel.title,
        novel.originalTitle,
        novel.cover,
        novel.source,
        novel.url,
        novel.summary,
        novel.author,
        novel.status,
        novel.chaptersCount,
        novel.lastRead?.chapterId,
        novel.lastRead?.progress,
        novel.lastRead?.timestamp,
        novel.lastUpdated,
        novel.dateAdded
      );
      
      // Insert genres if any
      if (novel.genres) {
        novel.genres.forEach(genre => {
          insertGenre.run(novel.id, genre);
        });
      }
    });
    
    // Insert chapters
    const insertChapter = db.prepare(`
      INSERT INTO chapters (
        id, novelId, number, title, originalTitle, 
        content, dateTranslated, wordCount
      )
      VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `);
    
    mockChapters.forEach(chapter => {
      insertChapter.run(
        chapter.id,
        chapter.novelId,
        chapter.number,
        chapter.title,
        chapter.originalTitle,
        chapter.content,
        chapter.dateTranslated,
        chapter.wordCount
      );
    });
    
    console.log('Database initialized with mock data');
  }
  
  return db;
};

// Database queries for novels
export const getAllNovels = (): Novel[] => {
  const db = getDb();
  const novels = db.prepare(`SELECT * FROM novels`).all() as any[];
  
  return novels.map(novel => {
    // Get genres for this novel
    const genres = db.prepare(`
      SELECT genre FROM novel_genres WHERE novelId = ?
    `).all(novel.id).map((row: any) => row.genre);
    
    // Transform database row to Novel type
    return {
      id: novel.id,
      title: novel.title,
      originalTitle: novel.originalTitle || undefined,
      cover: novel.cover || undefined,
      source: novel.source,
      url: novel.url,
      summary: novel.summary,
      author: novel.author || undefined,
      status: novel.status as 'Ongoing' | 'Completed' | 'Unknown' | undefined,
      genres: genres.length > 0 ? genres : undefined,
      chaptersCount: novel.chaptersCount,
      lastRead: novel.lastRead_chapterId ? {
        chapterId: novel.lastRead_chapterId,
        progress: novel.lastRead_progress,
        timestamp: novel.lastRead_timestamp
      } : undefined,
      lastUpdated: novel.lastUpdated,
      dateAdded: novel.dateAdded
    };
  });
};

export const getNovelById = (id: string): Novel | undefined => {
  const db = getDb();
  const novel = db.prepare(`SELECT * FROM novels WHERE id = ?`).get(id) as any;
  
  if (!novel) return undefined;
  
  // Get genres for this novel
  const genres = db.prepare(`
    SELECT genre FROM novel_genres WHERE novelId = ?
  `).all(id).map((row: any) => row.genre);
  
  // Transform database row to Novel type
  return {
    id: novel.id,
    title: novel.title,
    originalTitle: novel.originalTitle || undefined,
    cover: novel.cover || undefined,
    source: novel.source,
    url: novel.url,
    summary: novel.summary,
    author: novel.author || undefined,
    status: novel.status as 'Ongoing' | 'Completed' | 'Unknown' | undefined,
    genres: genres.length > 0 ? genres : undefined,
    chaptersCount: novel.chaptersCount,
    lastRead: novel.lastRead_chapterId ? {
      chapterId: novel.lastRead_chapterId,
      progress: novel.lastRead_progress,
      timestamp: novel.lastRead_timestamp
    } : undefined,
    lastUpdated: novel.lastUpdated,
    dateAdded: novel.dateAdded
  };
};

// Database queries for chapters
export const getChaptersForNovel = (novelId: string): Chapter[] => {
  const db = getDb();
  const chapters = db.prepare(`
    SELECT * FROM chapters 
    WHERE novelId = ?
    ORDER BY number ASC
  `).all(novelId) as any[];
  
  return chapters.map(chapter => ({
    id: chapter.id,
    novelId: chapter.novelId,
    number: chapter.number,
    title: chapter.title,
    originalTitle: chapter.originalTitle || undefined,
    content: chapter.content,
    dateTranslated: chapter.dateTranslated,
    wordCount: chapter.wordCount || undefined
  }));
};

export const getChapter = (novelId: string, chapterNumber: number): Chapter | undefined => {
  const db = getDb();
  const chapter = db.prepare(`
    SELECT * FROM chapters 
    WHERE novelId = ? AND number = ?
  `).get(novelId, chapterNumber) as any;
  
  if (!chapter) return undefined;
  
  return {
    id: chapter.id,
    novelId: chapter.novelId,
    number: chapter.number,
    title: chapter.title,
    originalTitle: chapter.originalTitle || undefined,
    content: chapter.content,
    dateTranslated: chapter.dateTranslated,
    wordCount: chapter.wordCount || undefined
  };
};

// Get all source sites
export const getAllSourceSites = (): SourceSite[] => {
  const db = getDb();
  const sites = db.prepare(`SELECT * FROM sources`).all() as any[];
  
  return sites.map(site => ({
    id: site.id,
    name: site.name,
    url: site.url,
    language: site.language as 'Chinese' | 'Korean' | 'Japanese' | 'Other',
    icon: site.icon || undefined
  }));
};
