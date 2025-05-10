package screenshots

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

type osProvider interface {
	getUserConfigDir() (string, error)
	mkdirAll(path string, perm os.FileMode) error
}

type realOsProvider struct {}

func (r *realOsProvider) getUserConfigDir() (string, error) {
	return os.UserConfigDir()
}

func (r *realOsProvider) mkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

type bleveProvider interface {
	Open(indexPath string) (indexer, error)
	New(path string, mapping mapping.IndexMapping) (indexer, error)
}

type realBleveProvider struct {}

func (r *realBleveProvider) Open(indexPath string) (indexer, error) {
    return bleve.Open(indexPath)
}

func (r *realBleveProvider) New(path string, mapping mapping.IndexMapping) (indexer, error) {
    return bleve.New(path, mapping)
}

type indexer interface {
	Index(id string, data interface{}) error
    Search(req *bleve.SearchRequest) (*bleve.SearchResult, error)
    Close() error
}

type IndexerProvider interface {
	Open() error
	Close() error
	Index(path string, doc *ScreenshotDoc) error
	Search(searchRequest *bleve.SearchRequest) (*bleve.SearchResult, error)
	GetIndexPath() (string, error)
}
type Indexer struct {
	appName string
	blevePath string
	idx indexer
	o osProvider
	b bleveProvider
}

func NewIndexer() *Indexer {
	return &Indexer{
		appName: "Glimpse",
		blevePath: "screenshots.bleve",

		o: &realOsProvider{},
		b: &realBleveProvider{},
	}
}

func (i *Indexer) GetIndexPath() (string, error) {
	userDataDir, err := i.o.getUserConfigDir()
	if err != nil {
		return "", err
	}

	appDataDir := filepath.Join(userDataDir, i.appName)
	if err := i.o.mkdirAll(appDataDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appDataDir, i.blevePath), nil
}

func (i *Indexer) Open() error {
	if i.idx != nil {
		return nil
	}
	
	indexPath, err := i.GetIndexPath()
	if err != nil {
		return fmt.Errorf("failed to create directory for index: %v", err)
	}

	idx, err := i.b.Open(indexPath)
	if err != nil {
		if err == bleve.ErrorIndexPathDoesNotExist {
			mapping := bleve.NewIndexMapping()
			idx, err = i.b.New(indexPath, mapping)
			if err != nil {
				return err
			}
		} else {
			return err
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
		return nil, err
	}

	return searchResult, nil
}

