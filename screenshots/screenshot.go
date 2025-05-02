package screenshots

import (
	"fmt"
	"path/filepath"
	"sync"
)

type Service interface {
	ScanAndIndex() error
}

type ScreenshotService struct {
	Dir DirProvider
	OCR OCRProvider
}

var supportedImageExts = map[string]struct{}{
    ".png":  {},
    ".jpg":  {},
    ".jpeg": {},
    ".gif":  {},
    ".bmp":  {},
    ".tif":  {},
    ".tiff": {},
    ".webp": {},
    ".svg":  {},
    ".heic": {},
    ".heif": {},
    ".avif": {},
}

func NewScreenshotService(d DirProvider, o OCRProvider) Service {
	return &ScreenshotService{
		Dir: d,
		OCR: o,
	}
}

func (s *ScreenshotService) ScanAndIndex() error {
	homeDir, err := s.Dir.GetHomeDir()
	if err != nil {
		return fmt.Errorf("error getting homedir: %w", err)
	}

	entries, err := s.Dir.ReadDir(homeDir)
	if err != nil {
		return fmt.Errorf("error reading screenshots dir: %w", err)
	}

	var wg sync.WaitGroup
	resultChan := make(chan string, len(entries))
	// errChan := make(chan error)

	for _, entry := range entries {
		fullPath := filepath.Join(filepath.Join(homeDir, "Desktop"), entry.Name())

		ext := filepath.Ext(fullPath)

		if _, ok := supportedImageExts[ext]; !ok {
			continue
		}

		wg.Add(1)
		go func(fullPath string){
			defer wg.Done()

			text, err := s.OCR.ExtractText(fullPath)
			if err != nil {
				return
			}

			resultChan <- text

			// select {
			// case resultChan <- text:
			// case errChan <- err:
			// }
		}(fullPath)
	}

	wg.Wait()
	close(resultChan)

	for r := range resultChan {
		fmt.Println(r)
	}

	return nil
}