# Arcane Translator

A modern web application for translating novels using AI, featuring a responsive React frontend and a high-performance Go backend.

## Features

- 🌐 **Web Scraping**: Extract content from various novel websites
- 🤖 **AI Translation**: Powered by Google Gemini for high-quality translations
- 📚 **Digital Library**: Organize and manage your novel collection
- 📖 **Reading Experience**: Chapter-by-chapter reading with progress tracking
- 🔍 **Search & Filter**: Find novels by title, author, genre, or language
- 🎨 **Modern UI**: Clean, responsive interface built with React and Tailwind CSS
- 💾 **Local Storage**: SQLite database for storing novels and reading progress

## Tech Stack

### Frontend
- **React 18** with TypeScript
- **Vite** for build tooling and development server
- **Tailwind CSS** for styling
- **shadcn/ui** for UI components
- **React Router** for navigation
- **TanStack Query** for data fetching

### Backend
- **Go 1.24+** with modern concurrency patterns
- **SQLite** database with better-sqlite3
- **Colly** for web scraping
- **Google Gemini AI** for translations

## Prerequisites

Before setting up the project, ensure you have the following installed:

- **Node.js 18.0+** and npm
- **Go 1.20+** (project uses Go 1.24)
- **Git**

### Verify Prerequisites

```bash
# Check Node.js version
node --version  # Should be 18.0 or higher

# Check npm version
npm --version

# Check Go version
go version  # Should be 1.20 or higher

# Check Git
git --version
```

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/DebanganThakuria/arcane-translator.git
cd arcane-translator
```

### 2. Setup Using Makefile (Recommended)

The project includes a Makefile with convenient commands:

```bash
# Install all dependencies (frontend + backend)
make install

# Build the backend
make build-backend
```

### 3. Manual Setup (Alternative)

If you prefer manual setup:

#### Frontend Dependencies
```bash
# Install React dependencies
npm install
```

#### Backend Dependencies
```bash
# Install Go dependencies
cd backend
go mod tidy
cd ..
```

#### Build Backend
```bash
# Build the Go server
cd backend
go build -o ../bin/server main.go
cd ..
```

## Running the Application

### Development Mode

For development with hot reloading:

#### Option 1: Using Makefile
```bash
# Start frontend development server
make start-frontend
```

#### Option 2: Manual Commands
```bash
# Start frontend (Terminal 1)
npm run dev

# Start backend (Terminal 2)
cd backend
go run main.go
```

### Production Mode

#### 1. Build Frontend
```bash
# Build optimized frontend
npm run build
```

#### 2. Start Backend Server
```bash
# Run the built server
./bin/server
```

The application will be available at:
- **Backend API**: http://localhost:8088
- **Frontend (dev)**: http://localhost:5173
- **Frontend (production)**: Served by backend at http://localhost:8088

## Project Structure

```
arcane-translator/
├── backend/                 # Go backend
│   ├── main.go             # Server entry point
│   ├── handler/            # HTTP handlers
│   ├── models/             # Data models
│   ├── provider/           # External service providers
│   │   ├── gemini/         # Google Gemini AI integration
│   │   ├── sources/        # Novel source scrapers
│   │   └── webscraper/     # Web scraping utilities
│   ├── repo/               # Database layer
│   ├── service/            # Business logic
│   └── utils/              # Utilities
├── src/                    # React frontend
│   ├── components/         # Reusable components
│   ├── pages/              # Page components
│   ├── services/           # Frontend services
│   ├── types/              # TypeScript types
│   └── utils/              # Frontend utilities
├── public/                 # Static assets
├── data/                   # Database files
└── bin/                    # Built binaries
```

## Configuration

### Environment Variables

The application may require certain environment variables for full functionality:

```bash
# Example .env file (create in project root if needed)
GEMINI_API_KEY=your_gemini_api_key_here
```

### Database

The application uses SQLite for data storage. The database file is automatically created in the `data/` directory when the backend starts.

## Available Commands

### Makefile Commands
```bash
make help           # Show available commands
make install        # Install all dependencies
make build-backend  # Build Go backend
make start-frontend # Start frontend development server
make stop          # Stop all running processes
make clean         # Clean dependencies and build files
```

### NPM Scripts
```bash
npm run dev        # Start Vite development server
npm run build      # Build for production
npm run preview    # Preview production build
npm run lint       # Run ESLint
```

## Usage

1. **Start the Application**: Follow the running instructions above
2. **Add Novels**: Use the "Add Novel" feature to import novels from supported websites
3. **Browse Library**: View your novel collection in the library
4. **Read & Translate**: Open novels to read and translate chapters
5. **Search**: Use the search functionality to find specific novels

## Supported Novel Sources

The application supports scraping from various novel websites. The scraping system includes:

- **Multi-source support**: Syosetu, Shuba Novel, and other platforms
- **Intelligent parsing**: Extracts chapter content and metadata
- **Error handling**: Robust error handling for failed requests

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Check what's using port 8088
   lsof -i :8088
   # Kill the process if needed
   kill -9 <PID>
   ```

2. **Go Dependencies Issues**
   ```bash
   cd backend
   go clean -modcache
   go mod download
   ```

3. **Frontend Build Issues**
   ```bash
   rm -rf node_modules package-lock.json
   npm install
   ```

4. **Database Permission Issues**
   ```bash
   # Ensure data directory exists and is writable
   mkdir -p data
   chmod 755 data
   ```

## Development

### Frontend Development
- Uses Vite for fast development with hot module replacement
- TypeScript for type safety
- ESLint for code linting
- Tailwind CSS for styling

### Backend Development
- Structured Go application with clean architecture
- RESTful API design
- SQLite for data persistence
- Comprehensive error handling

### Adding New Features
1. Backend: Add handlers in `backend/handler/`
2. Frontend: Add components in `src/components/` or pages in `src/pages/`
3. Database: Update models in `backend/models/`

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue in the GitHub repository or contact the maintainers.

---

**Happy Translating! 📚✨**
