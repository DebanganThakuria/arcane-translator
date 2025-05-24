
import { GoogleGenerativeAI } from '@google/generative-ai';

// Initialize the Google Generative AI with API key
let apiKey = '';

interface NovelDetails {
  novel_title_translated: string;
  novel_summary_translated: string;
  novel_author_name_translated?: string;
  possible_novel_genres?: string[];
  number_of_chapters?: number;
  chapter_links?: string[];
  status?: 'Ongoing' | 'Completed' | 'Unknown';
}

interface TranslatedChapter {
  translated_chapter_title: string;
  translated_chapter_contents: string;
  possible_new_genres?: string[];
}

// Configure the Gemini AI
const configureGemini = (key: string) => {
  apiKey = key;
};

const getGeminiClient = () => {
  if (!apiKey) {
    throw new Error('Gemini API key not configured');
  }
  return new GoogleGenerativeAI(apiKey);
};

// Extract and translate novel details from a webpage URL
// This will be moved to Python backend
export const extractNovelDetails = async (url: string, sourceSite: string): Promise<NovelDetails> => {
  try {
    // For demo purposes, we'll return mock data
    // In production, this will call the Python backend API
    return {
      novel_title_translated: "The Legendary Mechanic",
      novel_summary_translated: "As a space fleet mechanic, Han Xiao was one of millions working in the mechanical industry after Earth's apocalypse. A headset accident transported him into the body of an NPC he played in a game. With his foresight and natural-born calculating abilities, he manages to save a crucial character and changes the history of the universe's greatest war.",
      novel_author_name_translated: "Qi Peijia",
      possible_novel_genres: ["Sci-Fi", "Game", "Adventure", "Action"],
      number_of_chapters: 1465,
      status: "Completed"
    };
  } catch (error) {
    console.error('Error extracting novel details:', error);
    throw error;
  }
};

// Translate a chapter - will be moved to Python backend
export const translateChapter = async (novelId: string, chapterNumber: number, chapterUrl?: string, originalContent?: string): Promise<TranslatedChapter> => {
  try {
    // Mock response - will be replaced by Python backend API call
    return {
      translated_chapter_title: `Chapter ${chapterNumber}: New Beginning`,
      translated_chapter_contents: `<p>Han Xiao opened his eyes, feeling disoriented as he took in his surroundings...</p>`,
      possible_new_genres: ["Isekai", "LitRPG", "Sci-Fi", "Adventure"]
    };
  } catch (error) {
    console.error('Error translating chapter:', error);
    throw error;
  }
};

// Set the API Key for Gemini
export const setGeminiApiKey = (key: string) => {
  configureGemini(key);
};

// Check if the API key is configured
export const isGeminiConfigured = () => {
  return !!apiKey;
};

export default {
  extractNovelDetails,
  translateChapter,
  setGeminiApiKey,
  isGeminiConfigured
};
