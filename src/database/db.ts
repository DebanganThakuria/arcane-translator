import { Novel, Chapter, SourceSite } from '../types/novel';
import { mockNovels, mockChapters, mockSites } from '../data/mockData';

// Check if we're in a Node.js environment or browser
const isNode = typeof globalThis !== 'undefined' && globalThis.process && globalThis.process.versions && globalThis.process.versions.node;

// This will only be used in Node.js environment
let Database: any;
let fs: any;
let path: any;
let DB_PATH: string;

// Only load Node.js specific modules on the server side
if (isNode) {
  // Use top-level async IIFE (Immediately Invoked Function Expression) to load modules
  (async () => {
    try {
      console.log('Loading Node.js modules...');
      
      // Try to load modules using dynamic imports with explicit node: protocol
      let fsModule, pathModule, dbModule;
      
      try {
        // Use node: protocol which is more reliable in ESM context
        fsModule = await import('node:fs');
        console.log('Successfully loaded fs module');
      } catch (e) {
        console.error('Failed to load fs module:', e);
        fsModule = null;
      }
      
      try {
        pathModule = await import('node:path');
        console.log('Successfully loaded path module');
      } catch (e) {
        console.error('Failed to load path module:', e);
        pathModule = null;
      }
      
      // Assign the modules
      fs = fsModule;
      path = pathModule;
      
      if (!fs || !path) {
        throw new Error('Could not load fs or path modules');
      }
      
      // Set database path
      DB_PATH = path.join(process.cwd(), 'novels.db');
      console.log('Database path set to:', DB_PATH);
      
      // Try to load better-sqlite3
      try {
        dbModule = await import('better-sqlite3');
        console.log('Successfully loaded better-sqlite3 module');
      } catch (e) {
        console.error('Failed to load better-sqlite3 module:', e);
        dbModule = null;
      }
      
      // Assign the database module (handle both default and named exports)
      Database = dbModule?.default || dbModule;
      
      if (!Database) {
        throw new Error('Could not load better-sqlite3 module');
      }
      
      console.log('All Node.js modules loaded successfully. Initializing database...');
      
      // Initialize database now that modules are loaded
      initializeDatabase();
      
    } catch (error) {
      console.error('Failed to load required modules:', error);
      fs = null;
      path = null;
      Database = null;
      DB_PATH = '';
    }
  })().catch(error => {
    console.error('Error in async module loading:', error);
  });
}

// Database connection function that works in both environments
export const getDb = () => {
  if (isNode) {
    try {
      const db = new Database(DB_PATH);
      db.pragma('journal_mode = WAL');
      return db;
    } catch (error) {
      console.error('Failed to connect to SQLite database:', error);
      return null;
    }
  }
  return null; // Return null in browser environment
};

// Initialize the database schema and seed with mock data
export const initializeDatabase = () => {
  if (!isNode) {
    // console.log('Not in Node.js environment, skipping database initialization.');
    return null; 
  }

  // Add checks for fs, path, and Database being successfully loaded
  if (!fs || !path || !Database || !DB_PATH) {
    console.error('Essential Node.js modules (fs, path, better-sqlite3) or DB_PATH are not available. Cannot initialize database.');
    return null;
  }

  try {
    // Check if database already exists
    const dbExists = fs.existsSync(DB_PATH); // fs should now be defined if loaded
    
    const db = getDb(); // getDb also relies on Database and DB_PATH
    if (!db) {
      console.error('Failed to get DB connection in initializeDatabase.');
      return null;
    }
    
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
      
      mockChapters.forEach((chapter) => {
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
  } catch (error) {
    console.error('Failed to initialize database:', error);
    return null;
  }
};

// Database queries for novels - these work in both environments
export const getAllNovels = (): Novel[] => {
  if (isNode) {
    try {
      const db = getDb();
      if (!db) return mockNovels;
      
      const novels = db.prepare(`SELECT * FROM novels`).all() as any[];
      
      return novels.map((novel: any): Novel => {
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
    } catch (error) {
      console.error('Error getting novels from database:', error);
      return mockNovels;
    }
  }
  
  // In browser environment, return mock data
  return mockNovels;
};

export const getNovelById = (id: string): Novel | undefined => {
  if (isNode) {
    try {
      const db = getDb();
      if (!db) return mockNovels.find(novel => novel.id === id);
      
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
    } catch (error) {
      console.error('Error getting novel by ID from database:', error);
    }
  }
  
  // In browser environment, return mock data
  return mockNovels.find(novel => novel.id === id);
};

// Database queries for chapters
export const getChaptersForNovel = (novelId: string): Chapter[] => {
  if (isNode) {
    try {
      const db = getDb();
      if (!db) return mockChapters.filter(chapter => chapter.novelId === novelId);
      
      const chapters = db.prepare(`
        SELECT * FROM chapters 
        WHERE novelId = ?
        ORDER BY number ASC
      `).all(novelId) as any[];
      
      return chapters.map((chapter: any): Chapter => ({
        id: chapter.id,
        novelId: chapter.novelId,
        number: chapter.number,
        title: chapter.title,
        originalTitle: chapter.originalTitle || undefined,
        content: chapter.content,
        dateTranslated: chapter.dateTranslated,
        wordCount: chapter.wordCount || undefined
      }));
    } catch (error) {
      console.error('Error getting chapters from database:', error);
    }
  }
  
  // In browser environment, return mock data
  return mockChapters.filter(chapter => chapter.novelId === novelId);
};

export const getChapter = (novelId: string, chapterNumber: number): Chapter | undefined => {
  if (isNode) {
    try {
      const db = getDb();
      if (!db) {
        return mockChapters.find(
          chapter => chapter.novelId === novelId && chapter.number === chapterNumber
        );
      }
      
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
    } catch (error) {
      console.error('Error getting chapter from database:', error);
    }
  }
  
  // In browser environment, return mock data
  return mockChapters.find(
    chapter => chapter.novelId === novelId && chapter.number === chapterNumber
  );
};

// Get all source sites
export const getAllSourceSites = (): SourceSite[] => {
  if (isNode) {
    try {
      const db = getDb();
      if (!db) return mockSites;
      
      const sites = db.prepare(`SELECT * FROM sources`).all() as any[];
      
      return sites.map((site: any): SourceSite => ({
        id: site.id,
        name: site.name,
        url: site.url,
        language: site.language as 'Chinese' | 'Korean' | 'Japanese' | 'Other',
        icon: site.icon || undefined
      }));
    } catch (error) {
      console.error('Error getting source sites from database:', error);
    }
  }
  
  // In browser environment, return mock data
  return mockSites;
};
