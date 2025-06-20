# Arcane Translator

A modern web application for translating novels using AI, featuring a responsive React frontend and a high-performance Go backend.

## Features

- AI-powered novel translation
- Clean, responsive user interface
- Fast and efficient backend processing
- Cross-platform compatibility
- Easy installation via Homebrew

## Tech Stack

- **Frontend**: React, TypeScript, Vite
- **UI**: shadcn-ui, Tailwind CSS
- **Backend**: Go
- **Build Tool**: npm, Go modules

## Prerequisites

- Node.js 18+ and npm
- Go 1.20+
- Git

## Installation

### Using Homebrew (Recommended)


```bash
brew tap DebanganThakuria/tap
brew install arcane-translator
```

### Manual Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/DebanganThakuria/arcane-translator.git
   cd arcane-translator
   ```

2. Install dependencies and build:
   ```bash
   # Install frontend dependencies
   npm install
   
   # Build frontend
   npm run build
   
   # Install Go dependencies and build backend
   cd backend
   go mod download
   go build -o arcane-translator-backend
   ```

## Usage

1. Start the application:
   ```bash
   # If installed via Homebrew
   arcane-translator
   
   # For manual installation
   ./arcane-translator
   ```

2. Open your browser and navigate to:
   ```
   http://localhost:8088
   ```

## Development

### Frontend Development

```bash
# Start development server
npm run dev
```

### Backend Development

```bash
cd backend
go run main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue in the GitHub repository.
