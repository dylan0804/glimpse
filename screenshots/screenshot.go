package screenshots

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	rake "github.com/afjoseph/RAKE.go"
)

type Service interface {
	ScanAndIndex(context context.Context) error
}

type ScreenshotService struct {
	Dir DirProvider
	OCR OCRProvider
	Indexer IndexerProvider
}

type ScreenshotDoc struct {
	Path string
	Tags []string
}

var supportedImageExts = map[string]struct{}{
    ".png":  {},
    ".jpg":  {},
    ".jpeg": {},
    ".gif":  {},
    ".svg":  {},
}

func NewScreenshotService(d DirProvider, o OCRProvider, i IndexerProvider) Service {
	return &ScreenshotService{
		Dir: d,
		OCR: o,
		Indexer: i,
	}
}

func (s *ScreenshotService) ScanAndIndex(ctx context.Context) error {
	homeDir, err := s.Dir.GetHomeDir()
	if err != nil {
		return fmt.Errorf("error getting homedir: %v", err)
	}

	entries, err := s.Dir.ReadDir(homeDir)
	if err != nil {
		return fmt.Errorf("error reading screenshots dir: %v", err)
	}

	var wg sync.WaitGroup
	resultChan := make(chan ScreenshotDoc, len(entries))
	errChan := make(chan error)

	err = s.OCR.WriteOCRHelper()
	if err != nil {
		return fmt.Errorf("error reading binary for OR: %v", err)
	}

	err = s.Indexer.Open()
	if err != nil {
		return fmt.Errorf("error opening indexer: %v", err)
	}
	defer s.Indexer.Close()

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return fmt.Errorf("cancelling operation: %v", err)
		default:
		}

		fullPath := filepath.Join(filepath.Join(homeDir, "image-folder"), entry.Name())
		fmt.Println(fullPath)

		ext := filepath.Ext(fullPath)

		if _, ok := supportedImageExts[ext]; !ok {
			continue
		}

		wg.Add(1)
		go func(fullPath string){
			defer wg.Done()

			text, err := s.OCR.ExtractText(fullPath)
			if err != nil {
				errChan <- fmt.Errorf("error extracting text %v", err)
				return
			}

			candidates := rake.RunRake(text)

			// get the first 10 tags (sorted by how relevant it is)
			tags := make([]string, 10)
			for i := range candidates[:10] {
				tags[i] = candidates[i].Key
			}
			
			doc := ScreenshotDoc{
				Path: fullPath,
				Tags: tags,
			}

			err = s.Indexer.Index(doc.Path, &doc)
			if err != nil {
				errChan <- fmt.Errorf("error indexing image: %v", err)
				return
			}

			resultChan <- doc
		}(fullPath)
	}

	wg.Wait()
	close(resultChan)

	for r := range resultChan {
		fmt.Println(r)
	}

	return nil
}