/**
 * Not every novel has a cover, and there is no image source to fall back on.
 * Rather than showing the same grey book icon on every gap in the shelf, derive
 * a stable tint and a glyph from the novel itself, so a missing cover is still
 * recognisable and the shelf keeps its rhythm.
 */

const HUES = [
	{ h: 8, s: 34 }, // cinnabar
	{ h: 28, s: 30 }, // amber
	{ h: 96, s: 20 }, // moss
	{ h: 190, s: 26 }, // celadon
	{ h: 214, s: 28 }, // indigo
	{ h: 268, s: 22 }, // iris
	{ h: 340, s: 24 } // plum
];

function hash(input: string): number {
	let value = 2166136261;
	for (let i = 0; i < input.length; i++) {
		value ^= input.charCodeAt(i);
		value = Math.imul(value, 16777619);
	}
	return Math.abs(value);
}

export interface FallbackCover {
	glyph: string;
	background: string;
	ink: string;
}

export function fallbackCover(novel: { id: string; title: string }): FallbackCover {
	const seed = hash(novel.id || novel.title);
	const { h, s } = HUES[seed % HUES.length];
	const tilt = 130 + (seed % 60);

	// Initials from the translated title, so the whole interface stays in one
	// script. Two letters carry more identity than one at cover size.
	const words = novel.title
		.trim()
		.split(/\s+/)
		.filter((word) => /^[a-z]/i.test(word));

	const glyph =
		words.length >= 2
			? (words[0][0] + words[1][0]).toUpperCase()
			: (words[0]?.slice(0, 2) || novel.title.slice(0, 2) || '?').toUpperCase();

	return {
		glyph,
		background: `linear-gradient(${tilt}deg, hsl(${h} ${s}% 34%), hsl(${(h + 24) % 360} ${s - 6}% 20%))`,
		ink: `hsl(${h} ${Math.max(s - 14, 8)}% 88%)`
	};
}
