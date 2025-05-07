
import { Plugin } from 'vite';
import { initializeDatabase } from './db';

export function sqlitePlugin(): Plugin {
  return {
    name: 'vite-plugin-sqlite',
    configureServer() {
      // Initialize database on server start
      // This code runs on the server, not in the browser
      try {
        initializeDatabase();
        console.log('SQLite database initialized');
      } catch (error) {
        console.error('Failed to initialize SQLite database:', error);
      }
    },
    config: () => ({
      optimizeDeps: {
        exclude: ['better-sqlite3']
      }
    })
  };
}
