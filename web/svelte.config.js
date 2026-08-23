import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		// Single-page app: every route falls back to index.html and is resolved
		// client-side. The Go API on :8088 is the only backend.
		adapter: adapter({ fallback: 'index.html', strict: false }),
		alias: {
			$components: 'src/lib/components',
			$api: 'src/lib/api',
			$stores: 'src/lib/stores',
			$utils: 'src/lib/utils'
		}
	}
};

export default config;
