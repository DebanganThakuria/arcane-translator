
import { GoogleGenerativeAI } from '@google/generative-ai';

// Initialize the Google Generative AI with API key
// In a production environment, this should be stored in environment variables
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
export const extractNovelDetails = async (url: string, sourceSite: string): Promise<NovelDetails> => {
  try {
    // In a real implementation, we would fetch and parse the webpage
    // For demo purposes, we'll simulate this with a direct call to Gemini
    const genAI = getGeminiClient();
    const model = genAI.getGenerativeModel({ model: "gemini-1.5-pro" });
    
    const prompt = `
      You are a professional translator for novels. 
      Please extract and translate information from this ${sourceSite} novel page: ${url}
      
      Return ONLY a valid JSON object with this exact structure:
      {
        "novel_title_translated": "Translated title in English",
        "novel_summary_translated": "Translated summary in English",
        "novel_author_name_translated": "Translated author name if available",
        "possible_novel_genres": ["Genre1", "Genre2", ...],
        "number_of_chapters": estimated number of chapters as integer,
        "status": "Ongoing" or "Completed" or "Unknown"
      }
      
      Do not include any commentary, explanation, or preamble. Only return the JSON object.
    `;
    
    // For demo purposes, we'll return mock data instead of calling the API
    // In a real implementation, the following would be replaced with:
    // const result = await model.generateContent(prompt);
    // const responseText = result.response.text();
    // return JSON.parse(responseText);
    
    // Mock response for demonstration
    return {
      novel_title_translated: "The Legendary Mechanic",
      novel_summary_translated: "As a space fleet mechanic, Han Xiao was one of millions working in the mechanical industry after Earth's apocalypse. A headset accident transported him into the body of an NPC he played in a game. With his foresight and natural-born calculating abilities, he manages to save a crucial character and changes the history of the universe's greatest war. This intriguing sci-fi tale depicts the journey of Han Xiao: a lowly mechanic who, step by step, works his way to the greatest heights.",
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

// Translate a chapter by its URL or content
export const translateChapter = async (novelId: string, chapterNumber: number, chapterUrl?: string, originalContent?: string): Promise<TranslatedChapter> => {
  try {
    const genAI = getGeminiClient();
    const model = genAI.getGenerativeModel({ model: "gemini-1.5-pro" });
    
    let prompt;
    if (chapterUrl) {
      prompt = `
        You are a professional translator for novels.
        Please translate this chapter from ${chapterUrl} into fluent English.
        
        Return ONLY a valid JSON object with this structure:
        {
          "translated_chapter_title": "Translated chapter title",
          "translated_chapter_contents": "Full translated content with HTML paragraph tags",
          "possible_new_genres": ["Genre1", "Genre2"...]
        }
        
        Make sure the translation reads naturally in English. Preserve paragraph breaks with <p> tags.
        Only return the JSON object, no other text.
      `;
    } else if (originalContent) {
      prompt = `
        You are a professional translator for novels.
        Please translate this content into fluent English:
        
        ${originalContent.substring(0, 5000)} ${originalContent.length > 5000 ? '...(content truncated)' : ''}
        
        Return ONLY a valid JSON object with this structure:
        {
          "translated_chapter_title": "Translated chapter title",
          "translated_chapter_contents": "Full translated content with HTML paragraph tags",
          "possible_new_genres": ["Genre1", "Genre2"...]
        }
        
        Make sure the translation reads naturally in English. Preserve paragraph breaks with <p> tags.
        Only return the JSON object, no other text.
      `;
    } else {
      throw new Error('Either chapter URL or content must be provided');
    }
    
    // For demo purposes, we'll return mock data instead of calling the API
    // In a real implementation, the following would be replaced with:
    // const result = await model.generateContent(prompt);
    // const responseText = result.response.text();
    // return JSON.parse(responseText);
    
    // Mock response with different content based on chapter number
    return {
      translated_chapter_title: `Chapter ${chapterNumber}: New Beginning`,
      translated_chapter_contents: `<p>Han Xiao opened his eyes, feeling disoriented as he took in his surroundings. The last thing he remembered was putting on the virtual reality headset for a gaming session, but now he found himself in an unfamiliar room.</p><p>"System initialization complete. Welcome, user."</p><p>A mechanical voice echoed in his mind, startling him. He sat up abruptly, looking around for the source of the voice but finding nothing.</p><p>"What's going on?" he muttered to himself, trying to make sense of the situation.</p><p>As he stood up, memories that weren't his own began flooding into his consciousness. They belonged to a character named Black Star, an NPC mechanic in the game Galaxy he had been playing. Somehow, impossibly, he seemed to have transmigrated into the body of that character.</p><p>"This can't be happening," Han Xiao whispered, but the evidence was overwhelming. He recognized this room from the game—it was Black Star's workshop on Planet Aquamarine.</p><p>The system voice spoke again. "Character status: Level 1 Mechanic. Skills available: Basic Engineering, Weapon Crafting Level 1."</p><p>Han Xiao's mind raced as he tried to process what had happened to him. If he really had become Black Star, then he knew exactly what would happen next according to the game's plot. In less than three days, Planet Aquamarine would be invaded, and Black Star—a minor NPC—would die in the initial attack.</p><p>But Han Xiao had no intention of following that script. With his knowledge of future events, he now had a chance to change everything.</p><p>"I need to prepare," he said firmly, looking around the workshop with new determination. What had been a death sentence for Black Star could be the beginning of something extraordinary for Han Xiao.</p>`,
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
