// Pure client-side app: all data lives behind the Go API on :8088, and the
// build is a static shell that resolves every route in the browser.
export const ssr = false;
export const prerender = false;
export const trailingSlash = 'never';
