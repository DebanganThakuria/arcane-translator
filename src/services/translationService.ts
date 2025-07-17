// Standard library imports
import {Chapter, Novel} from '../types/novel';

// API base URL - will be replaced by actual backend calls
const API_BASE_URL = 'http://localhost:8088';

// Manual content extraction method
export const scrapeManually = async (url: string): Promise<string> => {
  return new Promise((resolve, reject) => {
    // Open the URL in a new tab
    const newTab = window.open(url, '_blank');
    
    // Create and show modal dialog
    const modal = document.createElement('div');
    modal.style.cssText = `
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0, 0, 0, 0.5);
      display: flex;
      justify-content: center;
      align-items: center;
      z-index: 10000;
      font-family: system-ui, -apple-system, sans-serif;
    `;
    
    modal.innerHTML = `
      <div style="
        background: white;
        padding: 24px;
        border-radius: 8px;
        width: 90%;
        max-width: 600px;
        max-height: 90vh;
        overflow-y: auto;
        box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
      ">
        <h2 style="margin: 0 0 16px 0; font-size: 20px; font-weight: 600;">Manual Content Extraction</h2>
        <p style="margin: 0 0 16px 0; color: #666; font-size: 14px;">
          The URL has been opened in a new tab. Please copy the page source and paste it below.
        </p>
        <div style="margin: 0 0 16px 0; padding: 8px; background: #f5f5f5; border-radius: 4px; font-family: monospace; font-size: 12px; word-break: break-all;">
          ${url}
        </div>
        <div style="margin: 0 0 8px 0; font-weight: 500; font-size: 14px;">Instructions:</div>
        <ol style="margin: 0 0 16px 0; padding-left: 20px; font-size: 14px; color: #666;">
          <li>Go to the opened tab</li>
          <li>Right-click and select "View Page Source" or press Ctrl+U</li>
          <li>Copy all content (Ctrl+A then Ctrl+C)</li>
          <li>Paste it in the text area below</li>
        </ol>
        <textarea id="contentInput" placeholder="Paste the webpage content here..." style="
          width: 100%;
          height: 200px;
          padding: 12px;
          border: 2px solid #e5e5e5;
          border-radius: 4px;
          font-family: monospace;
          font-size: 12px;
          resize: vertical;
          box-sizing: border-box;
        "></textarea>
        <div style="margin: 8px 0 16px 0; font-size: 12px; color: #666;">
          Characters: <span id="charCount">0</span> (minimum 50 required)
        </div>
        <div style="display: flex; gap: 12px; justify-content: flex-end;">
          <button id="cancelBtn" style="
            padding: 8px 16px;
            border: 2px solid #e5e5e5;
            background: white;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
          ">Cancel</button>
          <button id="submitBtn" style="
            padding: 8px 16px;
            border: none;
            background: #3b82f6;
            color: white;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
          " disabled>Extract Content</button>
        </div>
      </div>
    `;
    
    document.body.appendChild(modal);
    
    const textarea = modal.querySelector('#contentInput') as HTMLTextAreaElement;
    const submitBtn = modal.querySelector('#submitBtn') as HTMLButtonElement;
    const cancelBtn = modal.querySelector('#cancelBtn') as HTMLButtonElement;
    const charCount = modal.querySelector('#charCount') as HTMLSpanElement;
    
    // Focus on textarea
    textarea.focus();
    
    // Update character count and button state
    const updateState = () => {
      const length = textarea.value.trim().length;
      charCount.textContent = length.toString();
      submitBtn.disabled = length < 50;
      submitBtn.style.background = length < 50 ? '#9ca3af' : '#3b82f6';
      submitBtn.style.cursor = length < 50 ? 'not-allowed' : 'pointer';
    };
    
    textarea.addEventListener('input', updateState);
    
    // Handle submit
    submitBtn.addEventListener('click', () => {
      const content = textarea.value.trim();
      if (content.length >= 50) {
        // Close the opened tab
        if (newTab && !newTab.closed) {
          newTab.close();
        }
        
        // Remove modal
        document.body.removeChild(modal);
        
        // Resolve with content
        resolve(content);
      }
    });
    
    // Handle cancel
    cancelBtn.addEventListener('click', () => {
      // Close the opened tab
      if (newTab && !newTab.closed) {
        newTab.close();
      }
      
      // Remove modal
      document.body.removeChild(modal);
      
      // Reject
      reject(new Error('Manual extraction cancelled by user'));
    });
    
    // Handle escape key
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        cancelBtn.click();
      }
    };
    
    document.addEventListener('keydown', handleEscape);
    
    // Handle Ctrl+Enter to submit
    textarea.addEventListener('keydown', (e) => {
      if (e.ctrlKey && e.key === 'Enter' && !submitBtn.disabled) {
        submitBtn.click();
      }
    });
    
    // Clean up escape listener when modal is removed
    const originalRemove = () => {
      document.removeEventListener('keydown', handleEscape);
    };
    
    // Override the removal to clean up
    const cleanup = () => {
      originalRemove();
      if (newTab && !newTab.closed) {
        newTab.close();
      }
    };
    
    // Store cleanup function
    (modal as any).cleanup = cleanup;
  });
};

// Browser-based novel extraction
export const extractNovelDetailsBrowser = async (url: string, sourceSite: string): Promise<Novel> => {
  try {
    // First, validate the URL format
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      throw new Error('Invalid URL format. URL must start with http:// or https://');
    }

    // Scrape the content using manual extraction
    const htmlContent = await scrapeManually(url);
    
    const response = await fetch(`${API_BASE_URL}/novels/translate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ 
        url: url, 
        source: sourceSite,
        html_content: htmlContent 
      })
    });
    
    if (!response.ok) {
      let errorMessage = 'Failed to fetch novel details';
      
      try {
        const errorData = await response.json();
        
        // Handle specific error codes with user-friendly messages
        if (response.status === 404) {
          errorMessage = 'Novel not found. Please check if the URL is correct.';
        } else if (response.status === 403) {
          errorMessage = 'Access to this website is forbidden. The website may be blocking our requests.';
        } else if (response.status === 500) {
          errorMessage = errorData.detail || 'Server error occurred while extracting novel details.';
        } else {
          errorMessage = errorData.detail || errorMessage;
        }
      } catch (e) {
        // If parsing JSON fails, use a generic message
        if (response.status === 404) {
          errorMessage = 'Novel not found. Please check if the URL is correct.';
        } else if (response.status === 403) {
          errorMessage = 'Access to this website is forbidden. The website may be blocking our requests.';
        }
      }
      
      throw new Error(errorMessage);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error extracting novel details using browser:', error);
    throw error;
  }
};

// Browser-based first chapter translation
export const setFirstChapterUrlBrowser = async (
  novelId: string, 
  firstChapterUrl: string,
  htmlContent: string,
): Promise<Chapter> => {
  try {
    const response = await fetch(
      `${API_BASE_URL}/novels/translate/first_chapter`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
          novel_id: novelId, 
          chapter_url: firstChapterUrl,
          html_content: htmlContent 
        })
      }
    );
    
    if (!response.ok) {
      let errorMessage = 'Failed to set first chapter URL';
      
      try {
        const errorData = await response.json();
        errorMessage = errorData.detail || errorMessage;
      } catch (e) {
        // If parsing JSON fails, use the generic message
      }
      
      throw new Error(errorMessage);
    }
    
    return await response.json();
  } catch (error) {
    console.error(`Error setting first chapter URL for novel ID ${novelId} using browser:`, error);
    throw error;
  }
};

// Extract and translate novel details from a webpage URL
export const extractNovelDetails = async (url: string, sourceSite: string): Promise<Novel> => {
  try {
    // First, validate the URL format
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      throw new Error('Invalid URL format. URL must start with http:// or https://');
    }
    
    const response = await fetch(`${API_BASE_URL}/novels/translate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ url: url, source: sourceSite })
    });
    
    if (!response.ok) {
      let errorMessage = 'Failed to fetch novel details';
      
      try {
        const errorData = await response.json();
        
        // Handle specific error codes with user-friendly messages
        if (response.status === 404) {
          errorMessage = 'Novel not found. Please check if the URL is correct.';
        } else if (response.status === 403) {
          errorMessage = 'Access to this website is forbidden. The website may be blocking our requests.';
        } else if (response.status === 500) {
          errorMessage = errorData.detail || 'Server error occurred while extracting novel details.';
        } else {
          errorMessage = errorData.detail || errorMessage;
        }
      } catch (e) {
        // If parsing JSON fails, use a generic message
        if (response.status === 404) {
          errorMessage = 'Novel not found. Please check if the URL is correct.';
        } else if (response.status === 403) {
          errorMessage = 'Access to this website is forbidden. The website may be blocking our requests.';
        }
      }
      
      throw new Error(errorMessage);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error extracting novel details:', error);
    throw error;
  }
};

// Extract novel details with automatic fallback to browser-based scraping
export const extractNovelDetailsWithFallback = async (url: string, sourceSite: string): Promise<Novel> => {
  try {
    // Try normal scraping first
    return await extractNovelDetails(url, sourceSite);
  } catch (error) {
    console.log('Normal scraping failed, attempting manual extraction...');
    
    // If normal scraping fails, try manual extraction
    try {
      return await extractNovelDetailsBrowser(url, sourceSite);
    } catch (manualError) {
      console.error('Manual extraction also failed:', manualError);
      throw new Error(`Normal scraping failed. Manual extraction error: ${manualError.message}`);
    }
  }
};

// Browser-based chapter translation
export const translateChapterBrowser = async (
  novelId: string,
  chapterUrl: string
): Promise<Chapter> => {
  try {
    // Validate URL format
    if (!chapterUrl.startsWith('http://') && !chapterUrl.startsWith('https://')) {
      throw new Error('Invalid URL format. URL must start with http:// or https://');
    }

    // Scrape the content using manual extraction
    const htmlContent = await scrapeManually(chapterUrl);

    const response = await fetch(`${API_BASE_URL}/novels/translate/chapter`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ 
        novel_id: novelId,
        chapter_url: chapterUrl,
        html_content: htmlContent 
      })
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ detail: 'Unknown error' }));
      throw new Error(errorData.detail || 'Failed to translate chapter');
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error translating chapter using browser:', error);
    throw error;
  }
};

// Browser-based novel refresh
export const refreshNovelBrowser = async (
  novelId: string,
  novelUrl: string
): Promise<Novel> => {
  try {
    // Validate URL format
    if (!novelUrl.startsWith('http://') && !novelUrl.startsWith('https://')) {
      throw new Error('Invalid URL format. URL must start with http:// or https://');
    }

    // Scrape the content using manual extraction
    const htmlContent = await scrapeManually(novelUrl);

    const response = await fetch(`${API_BASE_URL}/novels/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        novel_id: novelId,
        html_content: htmlContent 
      })
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ detail: 'Unknown error' }));
      throw new Error(errorData.detail || 'Failed to refresh novel');
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error refreshing novel using browser:', error);
    throw error;
  }
};

// Translate a chapter
export const translateChapter = async (
  novelId: string,
): Promise<Chapter> => {
  try {
    const response = await fetch(`${API_BASE_URL}/novels/translate/chapter`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ novel_id: novelId})
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ detail: 'Unknown error' }));
      throw new Error(errorData.detail || 'Failed to translate chapter');
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error translating chapter:', error);
    throw error;
  }
};

// Translate chapter with automatic fallback to browser-based scraping
export const translateChapterWithFallback = async (
  novelId: string,
  chapterUrl?: string
): Promise<Chapter> => {
  try {
    // Try normal scraping first
    return await translateChapter(novelId);
  } catch (error) {
    console.log('Normal chapter scraping failed, attempting manual extraction...');
    
    // If normal scraping fails, try manual extraction
    if (!chapterUrl) {
      throw new Error('Chapter URL is required for manual extraction fallback');
    }
    
    try {
      return await translateChapterBrowser(novelId, chapterUrl);
    } catch (manualError) {
      console.error('Manual extraction also failed:', manualError);
      throw new Error(`Normal scraping failed. Manual extraction error: ${manualError.message}`);
    }
  }
};

// Refresh novel with automatic fallback to browser-based scraping
export const refreshNovelWithFallback = async (
  novelId: string,
  novelUrl: string
): Promise<Novel> => {
  try {
    // Try normal scraping first
    const response = await fetch(`${API_BASE_URL}/novels/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        novel_id: novelId,
      })
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ detail: 'Unknown error' }));
      throw new Error(errorData.detail || 'Failed to refresh novel');
    }
    
    return await response.json();
  } catch (error) {
    console.log('Normal novel refresh failed, attempting manual extraction...');
    
    // If normal scraping fails, try manual extraction
    try {
      return await refreshNovelBrowser(novelId, novelUrl);
    } catch (manualError) {
      console.error('Manual extraction also failed:', manualError);
      throw new Error(`Normal refresh failed. Manual extraction error: ${manualError.message}`);
    }
  }
};

// Set the first chapter URL for a novel
export const setFirstChapterUrl = async (
  novelId: string, 
  firstChapterUrl: string, 
  autoTranslate: boolean = true
): Promise<Chapter> => {
  try {
    // Validate URL format
    if (!firstChapterUrl.startsWith('http://') && !firstChapterUrl.startsWith('https://')) {
      throw new Error('Invalid URL format. URL must start with http:// or https://');
    }

    const response = await fetch(
      `${API_BASE_URL}/novels/translate/first_chapter`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ novel_id: novelId, chapter_url: firstChapterUrl })
      }
    );
    
    if (!response.ok) {
      let errorMessage = 'Failed to set first chapter URL';
      
      try {
        const errorData = await response.json();
        errorMessage = errorData.detail || errorMessage;
      } catch (e) {
        // If parsing JSON fails, use the generic message
      }
      
      throw new Error(errorMessage);
    }
    
    return await response.json();
  } catch (error) {
    console.error(`Error setting first chapter URL for novel ID ${novelId}:`, error);
    throw error;
  }
};

export default {
  extractNovelDetails,
  translateChapter,
  setFirstChapterUrl,
};
