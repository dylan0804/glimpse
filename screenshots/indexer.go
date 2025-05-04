package screenshots

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
)

const (
	path = "screenshots.bleve"
)

type IndexerProvider interface {
	Open() error
	Close() error
	Index(string, *ScreenshotDoc) error
}

type Indexer struct {
	idx bleve.Index
}

func NewIndexer() IndexerProvider {
	return &Indexer{}
}

func (i *Indexer) Open() error {
	idx, err := bleve.Open(path)
	if err != nil {
		if err == bleve.ErrorIndexPathDoesNotExist {
			mapping := bleve.NewIndexMapping()
			idx, err = bleve.New(path, mapping)
			if err != nil {
				return fmt.Errorf("failed to create index: %v", err)
			}
		}
	}

	i.idx = idx

	return nil
}

func (i *Indexer) Close() error {
	if i.Index != nil {
		return i.idx.Close()
	}
	return nil
}

func (i *Indexer) Index(path string, doc *ScreenshotDoc) error {
	return i.idx.Index(path, doc)
}

