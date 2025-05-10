# Glimpse

Ever spent ages scrolling through hundreds of screenshots trying to find that ONE image? "I know it had the word 'metrics' in it somewhere..." 🤔

Glimpse solves this headache! This magical tool extracts text from your screenshots, remembers what's important, and lets you search by the words that appeared in your images. No more endless scrolling - just type what you remember and find it in seconds!

## 🎯 Super Powers

- **X-Ray Vision**: Our OCR technology reads text from your screenshots like magic!
- **Mind Reading**: RAKE algorithm figures out what your screenshots are actually about
- **Time Travel**: Lightning-fast Bleve search finds that screenshot from three months ago in milliseconds

## 🧰 Under the Hood

- **Mighty Brain**: Powered by Go for that extra zoom
- **Eagle Eyes**: Tesseract OCR reads even the tiniest text
- **Smart Cookies**: RAKE algorithm knows what matters in your content
- **Supersonic Search**: Bleve search engine finds needles in haystacks

## 🚀 Get Started

### What You'll Need

- Go 1.18+
- Wails CLI
- Node.js and npm
- Tesseract OCR engine:
  - **macOS**: `brew install tesseract` 🍎

### The Spellbook

Clone this magical repository:
```bash
git clone https://github.com/yourusername/glimpse.git
cd glimpse
```

Gather your ingredients:
```bash
go mod tidy
```

Bring it to life:
```bash
wails dev
```

### Craft Your Potion

```bash
wails build
```

### Test Your Creation

```bash
go test -v ./screenshots
```

## 📜 License

MIT (Make It Tremendous!)
