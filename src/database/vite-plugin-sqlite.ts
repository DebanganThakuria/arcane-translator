
import { Plugin } from 'vite';
import { initializeDatabase } from './db';

export function sqlitePlugin(): Plugin {
  return {
    name: 'vite-plugin-sqlite',
    configureServer() {
      // Initialize database on server start
      initializeDatabase();
      console.log('SQLite database initialized');
    }
  };
}
