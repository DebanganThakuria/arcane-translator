import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		// Defaults to 8080 to match how the app has always been served, but PORT
		// wins so the dev server can be moved out of the way when it is taken.
		port: Number(process.env.PORT) || 8080,
		host: true
	}
});
