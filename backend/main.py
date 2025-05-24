
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional
import sqlite3
import os
from datetime import datetime
import requests
from bs4 import BeautifulSoup
import google.generativeai as genai
from dotenv import load_dotenv

load_dotenv()

app = FastAPI(title="Novel Translation API")

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:8080", "http://127.0.0.1:8080"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Configure Gemini AI
genai.configure(api_key=os.getenv("GEMINI_API_KEY"))

# Pydantic models
class Novel(BaseModel):
    id: str
    title: str
    originalTitle: Optional[str] = None
    cover: Optional[str] = None
    source: str
    url: str
    summary: str
    author: Optional[str] = None
    status: Optional[str] = None
    genres: Optional[List[str]] = None
    chaptersCount: int
    lastUpdated: int
    dateAdded: int

class Chapter(BaseModel):
    id: str
    novelId: str
    number: int
    title: str
    originalTitle: Optional[str] = None
    content: str
    dateTranslated: int
    wordCount: Optional[int] = None

class SourceSite(BaseModel):
    id: str
    name: str
    url: str
    language: str
    icon: Optional[str] = None

# Database setup
def init_db():
    conn = sqlite3.connect('novels.db')
    cursor = conn.cursor()
    
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS novels (
            id TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            originalTitle TEXT,
            cover TEXT,
            source TEXT NOT NULL,
            url TEXT NOT NULL,
            summary TEXT NOT NULL,
            author TEXT,
            status TEXT,
            chaptersCount INTEGER NOT NULL DEFAULT 0,
            lastUpdated INTEGER NOT NULL,
            dateAdded INTEGER NOT NULL
        )
    ''')
    
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS chapters (
            id TEXT PRIMARY KEY,
            novelId TEXT NOT NULL,
            number INTEGER NOT NULL,
            title TEXT NOT NULL,
            originalTitle TEXT,
            content TEXT NOT NULL,
            dateTranslated INTEGER NOT NULL,
            wordCount INTEGER,
            FOREIGN KEY (novelId) REFERENCES novels(id)
        )
    ''')
    
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS sources (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            url TEXT NOT NULL,
            language TEXT NOT NULL,
            icon TEXT
        )
    ''')
    
    conn.commit()
    conn.close()

# Initialize database on startup
init_db()

@app.get("/")
async def root():
    return {"message": "Novel Translation API"}

@app.get("/novels", response_model=List[Novel])
async def get_novels():
    conn = sqlite3.connect('novels.db')
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM novels")
    novels = cursor.fetchall()
    conn.close()
    
    return [
        Novel(
            id=novel[0],
            title=novel[1],
            originalTitle=novel[2],
            cover=novel[3],
            source=novel[4],
            url=novel[5],
            summary=novel[6],
            author=novel[7],
            status=novel[8],
            chaptersCount=novel[9],
            lastUpdated=novel[10],
            dateAdded=novel[11]
        ) for novel in novels
    ]

@app.get("/novels/{novel_id}", response_model=Novel)
async def get_novel(novel_id: str):
    conn = sqlite3.connect('novels.db')
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM novels WHERE id = ?", (novel_id,))
    novel = cursor.fetchone()
    conn.close()
    
    if not novel:
        raise HTTPException(status_code=404, detail="Novel not found")
    
    return Novel(
        id=novel[0],
        title=novel[1],
        originalTitle=novel[2],
        cover=novel[3],
        source=novel[4],
        url=novel[5],
        summary=novel[6],
        author=novel[7],
        status=novel[8],
        chaptersCount=novel[9],
        lastUpdated=novel[10],
        dateAdded=novel[11]
    )

@app.get("/novels/{novel_id}/chapters", response_model=List[Chapter])
async def get_chapters(novel_id: str):
    conn = sqlite3.connect('novels.db')
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM chapters WHERE novelId = ? ORDER BY number ASC", (novel_id,))
    chapters = cursor.fetchall()
    conn.close()
    
    return [
        Chapter(
            id=chapter[0],
            novelId=chapter[1],
            number=chapter[2],
            title=chapter[3],
            originalTitle=chapter[4],
            content=chapter[5],
            dateTranslated=chapter[6],
            wordCount=chapter[7]
        ) for chapter in chapters
    ]

@app.post("/novels/extract")
async def extract_novel_details(url: str, source: str):
    """Extract and translate novel details from a webpage URL"""
    try:
        # Fetch webpage content
        response = requests.get(url)
        soup = BeautifulSoup(response.content, 'html.parser')
        
        # Extract content (this would be customized per source site)
        title = soup.find('title').text if soup.find('title') else "Unknown Title"
        
        # Use Gemini to translate and extract details
        model = genai.GenerativeModel('gemini-1.5-pro')
        
        prompt = f"""
        You are a professional translator for novels. 
        Please extract and translate information from this {source} novel page content.
        
        Page title: {title}
        Page URL: {url}
        
        Return ONLY a valid JSON object with this exact structure:
        {{
            "novel_title_translated": "Translated title in English",
            "novel_summary_translated": "Translated summary in English",
            "novel_author_name_translated": "Translated author name if available",
            "possible_novel_genres": ["Genre1", "Genre2", ...],
            "number_of_chapters": estimated number of chapters as integer,
            "status": "Ongoing" or "Completed" or "Unknown"
        }}
        
        Do not include any commentary, explanation, or preamble. Only return the JSON object.
        """
        
        response = model.generate_content(prompt)
        
        # For demo, return mock data
        return {
            "novel_title_translated": "The Legendary Mechanic",
            "novel_summary_translated": "Epic sci-fi adventure about a mechanic in a game world.",
            "novel_author_name_translated": "Qi Peijia",
            "possible_novel_genres": ["Sci-Fi", "Game", "Adventure"],
            "number_of_chapters": 1465,
            "status": "Completed"
        }
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error extracting novel details: {str(e)}")

@app.post("/chapters/translate")
async def translate_chapter(novel_id: str, chapter_number: int, chapter_url: str = None, original_content: str = None):
    """Translate a chapter"""
    try:
        # For demo, return mock data
        return {
            "translated_chapter_title": f"Chapter {chapter_number}: New Beginning",
            "translated_chapter_contents": "<p>Translated chapter content would go here...</p>",
            "possible_new_genres": ["Action", "Adventure"]
        }
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error translating chapter: {str(e)}")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
