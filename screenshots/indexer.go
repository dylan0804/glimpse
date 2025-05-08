package screenshots

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
)
type IndexerProvider interface {
	Open() error
	Close() error
	Index(path string, doc *ScreenshotDoc) error
	Search(searchRequest *bleve.SearchRequest) (*bleve.SearchResult, error)
	GetIndexPath() (string, error)
}

type Indexer struct {
	idx bleve.Index
}

func NewIndexer() IndexerProvider {
	return &Indexer{}
}

func (i *Indexer) GetIndexPath() (string, error) {
	userDataDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDataDir := filepath.Join(userDataDir, "Glimpse")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appDataDir, "screenshots.bleve"), nil
}

func (i *Indexer) Open() error {
	if i.idx != nil {
		return nil
	}

	indexPath, err := i.GetIndexPath()
	if err != nil {
		return fmt.Errorf("failed to create directory for index: %v", err)
	}

	idx, err := bleve.Open(indexPath)
	if err != nil {
		if err == bleve.ErrorIndexPathDoesNotExist {
			mapping := bleve.NewIndexMapping()
			idx, err = bleve.New(indexPath, mapping)
			if err != nil {
				return fmt.Errorf("failed to create index: %v", err)
			}
		}
	}

	i.idx = idx

	return nil
}

func (i *Indexer) Close() error {
	return i.idx.Close()
}

func (i *Indexer) Index(path string, doc *ScreenshotDoc) error {
	return i.idx.Index(path, doc)
}

func (i *Indexer) Search(request *bleve.SearchRequest) (*bleve.SearchResult, error) {
	searchResult, err := i.idx.Search(request)
	if err != nil {
		return nil, fmt.Errorf("error when searching through index: %v", err)
	}

	return searchResult, nil
}

