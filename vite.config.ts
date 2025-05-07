
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';
import { sqlitePlugin } from './src/database/vite-plugin-sqlite';
import { componentTagger } from "lovable-tagger";

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => ({
  plugins: [
    react(),
    mode === 'development' && componentTagger(),
    sqlitePlugin()
  ].filter(Boolean),
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  optimizeDeps: {
    exclude: ['better-sqlite3'],
    include: ['@google/generative-ai']
  },
  server: {
    port: 8080,
    host: "::",
    watch: {
      ignored: ['**/novels.db', '**/novels.db-shm', '**/novels.db-wal']
    }
  }
}));
