
# Novel Translation Backend

Python FastAPI backend for the novel translation application.

## Setup

1. Install Python dependencies:
```bash
cd backend
pip install -r requirements.txt
```

2. Create a `.env` file with your Gemini API key:
```
GEMINI_API_KEY=your_api_key_here
```

3. Run the server:
```bash
python main.py
```

The API will be available at `http://localhost:8000`

## API Documentation

Once running, visit `http://localhost:8000/docs` for the interactive API documentation.

## Endpoints

- `GET /novels` - Get all novels
- `GET /novels/{id}` - Get a specific novel
- `GET /novels/{id}/chapters` - Get chapters for a novel
- `POST /novels/extract` - Extract novel details from URL
- `POST /chapters/translate` - Translate a chapter
