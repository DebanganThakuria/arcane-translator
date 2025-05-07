
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';
import { sqlitePlugin } from './src/database/vite-plugin-sqlite';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), sqlitePlugin()],
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
    watch: {
      ignored: ['**/novels.db', '**/novels.db-shm', '**/novels.db-wal']
    }
  }
});
