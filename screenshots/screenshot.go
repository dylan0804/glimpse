package screenshots

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	rake "github.com/afjoseph/RAKE.go"
	"github.com/blevesearch/bleve/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	b64 "encoding/base64"
)

type Service interface {
	ScanAndIndex() error
	Search(query string) error
	Shutdown()
}

type ScreenshotService struct {
	Dir DirProvider
	OCR OCRProvider
	Indexer IndexerProvider

	ctx context.Context
}

type ScreenshotDoc struct {
	Path string `json:"path"`
	Tags []string `json:"tags"`
	URL string `json:"url"`
}

var supportedImageExts = map[string]struct{}{
    ".png":  {},
    ".jpg":  {},
    ".jpeg": {},
    ".gif":  {},
    ".svg":  {},
}

func NewScreenshotService(d DirProvider, o OCRProvider, i IndexerProvider, ctx context.Context) Service {
	return &ScreenshotService{
		Dir: d,
		OCR: o,
		Indexer: i,
		ctx: ctx,
	}
}

func (s *ScreenshotService) ScanAndIndex() error {
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
	errChan := make(chan error, 100)

	err = s.OCR.WriteOCRHelper()
	if err != nil {
		return fmt.Errorf("error reading binary for OR: %v", err)
	}

	err = s.Indexer.Open()
	if err != nil {
		return fmt.Errorf("error opening indexer: %v", err)
	}

	go func(){
		for r := range resultChan {
			fmt.Println(r)
			runtime.EventsEmit(s.ctx, "result:found", r)
		}  
	}()

	go func(){
		for e := range errChan {
			fmt.Println(e)
		}  
	}()

	for _, entry := range entries {
		fullPath := filepath.Join(filepath.Join(homeDir, "Desktop"), entry.Name())

		ext := filepath.Ext(fullPath)

		if _, ok := supportedImageExts[ext]; !ok {
			continue
		}

		wg.Add(1)
		go func(fullPath string){
			defer wg.Done()

			// handle cancellation
			select {
			case <-s.ctx.Done():
				return
			default:
			}

			text, err := s.OCR.ExtractText(fullPath)
			if err != nil {
				errChan <- fmt.Errorf("error extracting text %v", err)
				return
			}
			if len(text) == 0 { // skip screenshots with no texts
				return
			}

			candidates := rake.RunRake(text)

			// get the first 20 tags (sorted by how relevant it is)
			tags := make([]string, 0, 20)
			
			for i := 0; i < cap(tags) && i < len(candidates); i++ {
				tags = append(tags, candidates[i].Key)
			}

			bytes, err := os.ReadFile(fullPath)
			if err != nil {
				errChan <- fmt.Errorf("error reading file: %v", err)
				return
			}
			
			doc := ScreenshotDoc{
				Path: fullPath,
				Tags: tags,
				URL: b64.StdEncoding.EncodeToString(bytes),
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
	close(errChan)

	return nil
}

func (s *ScreenshotService) Search(keyword string) error {
	err := s.Indexer.Open()
	if err != nil {
		return fmt.Errorf("error opening indexer: %v", err)
	}

	query := bleve.NewQueryStringQuery(keyword)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}

	searchResult, err := s.Indexer.Search(searchRequest)
	if err != nil {
		return err
	}

	for _, d := range searchResult.Hits {
		doc := ScreenshotDoc{
			Path: d.ID,
			URL: d.Fields["url"].(string),
		}

		runtime.EventsEmit(s.ctx, "search:found", doc)
	}

	return nil
}

func (s *ScreenshotService) Shutdown() {
	s.Indexer.Close()
}