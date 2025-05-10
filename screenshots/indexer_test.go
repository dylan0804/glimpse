package screenshots

import (
	"errors"
	"os"
	"testing"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search"
)

func NewMockIndexer(appName, blevePath string, o *mockOsProvider, b *mockBleveProvider, i *mockIndexer) *Indexer {
	return &Indexer{
		appName: appName,
		blevePath: blevePath,
		o: o,
		b: b,
		idx: i,
	}
}

type mockOsProvider struct {
	userConfigDir string

	userConfigDirErr error
	mkdirErr error
}

func (m *mockOsProvider) getUserConfigDir() (string, error) {
	if m.userConfigDirErr != nil {
		return "", m.userConfigDirErr
	}
	return m.userConfigDir, nil
}

func (m *mockOsProvider) mkdirAll(path string, perm os.FileMode) error {
	if m.mkdirErr != nil {
		return m.mkdirErr
	}
	return nil
}

type mockIndexer struct {
	indexError error
	searchError error
	closeError error

	searchFn func(req *bleve.SearchRequest) (*bleve.SearchResult, error)
}

func (m *mockIndexer) Index(id string, data interface{}) error {
	if m.indexError != nil {
		return m.indexError
	}
	return nil
}

func (m *mockIndexer) Search(req *bleve.SearchRequest) (*bleve.SearchResult, error) {
	if m.searchError != nil {
		return nil, m.searchError
	}
	return m.searchFn(req)
}

func (m *mockIndexer) Close() error {
	if m.closeError != nil {
		return m.closeError
	}
	return nil
}

type mockBleveProvider struct {
	openError error
	newError error
}

func (m *mockBleveProvider) Open(indexPath string) (indexer, error) {
	if m.openError != nil {
		return nil, m.openError
	}
	return &mockIndexer{}, nil
}

func (m *mockBleveProvider) New(path string, mapping mapping.IndexMapping) (indexer, error) {
	if m.newError != nil {
		return nil, m.newError
	}
	return &mockIndexer{}, nil
}


func TestGetIndexPath(t *testing.T){
	t.Run("Success", func(t *testing.T) {
		o := &mockOsProvider{
			userConfigDir: "user/config/dir",
		}

		i := NewMockIndexer("glimpse-test", "test.bleve", o, nil, nil)

		indexPath, err := i.GetIndexPath()
		if err != nil {
			t.Errorf("expected no error, got %v: ", err)
		}	
		if indexPath != "user/config/dir/glimpse-test/test.bleve" {
			t.Errorf("expected 'user/config/dir/glimpse-test/test.bleve', got: %s", indexPath)
		}
	})

	t.Run("UserConfigDir error", func(t *testing.T) {
		o := &mockOsProvider{
			userConfigDirErr: errors.New("user config dir error"),
		}
		i := NewMockIndexer("", "", o, nil, nil)

		_, err := i.GetIndexPath()
		if err == nil || err.Error() != "user config dir error" {
			t.Errorf("expected error 'user config dir error', got: %v", err)
		}
	})

	t.Run("Mkdir error", func(t *testing.T) {
		o := &mockOsProvider{
			mkdirErr: errors.New("mkdir error"),
		}
		i := NewMockIndexer("", "", o, nil, nil)

		_, err := i.GetIndexPath()
		if err == nil || err.Error() != "mkdir error" {
			t.Errorf("expected error 'mkdir error', got: %v", err)
		}
	})
}

func TestOpen(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		b := &mockBleveProvider{}
		o := &mockOsProvider{
			userConfigDir: "user/config/dir",
		}
		i := NewMockIndexer("glimpse-test", "test.bleve", o, b, nil)

		err := i.Open()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("PathDoesntExist error", func(t *testing.T) {
		b := &mockBleveProvider{
			openError: bleve.ErrorIndexPathDoesNotExist,
		}
		o := &mockOsProvider{
			userConfigDir: "user/config/dir",
		}
		i := NewMockIndexer("glimpse-test", "test.bleve", o, b, nil)
		
		err := i.Open()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("BleveNew error", func(t *testing.T) {
		b := &mockBleveProvider{
			openError: bleve.ErrorIndexPathDoesNotExist,
			newError: errors.New("error creating idx with mapping"),
		}
		o := &mockOsProvider{
			userConfigDir: "user/config/dir",
		}
		i := NewMockIndexer("glimpse-test", "test.bleve", o, b, nil)

		err := i.Open()
		if err == nil || err.Error() != "error creating idx with mapping" {
			t.Errorf("expected error 'error creating idx with mapping', got: %v", err)
		}
	})
}

func TestSearch(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockIdx := &mockIndexer{
			searchFn: func(req *bleve.SearchRequest) (*bleve.SearchResult, error) {
				return &bleve.SearchResult{
					Total: 1,
					Hits: search.DocumentMatchCollection{
						&search.DocumentMatch{ID: "test-id"},
					},
				}, nil
			},
		}

		i := NewMockIndexer("", "", nil, nil, mockIdx)

		searchResult, err := i.Search(nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if searchResult.Hits[0].ID != "test-id" {
			t.Errorf("expected 'test-id', got %s", searchResult.Hits[0].ID)
		}
	})

	t.Run("Search error", func(t *testing.T) {
		mockIdx := &mockIndexer{
			searchError: errors.New("search error"),
		}

		i := NewMockIndexer("", "", nil, nil, mockIdx)

		_, err := i.Search(nil)
		if err == nil || err.Error() != "search error" {
			t.Errorf("expected 'search error', got: %v", err)
		}
	})
}

func TestIndex(t *testing.T){
	t.Run("Success", func(t *testing.T) {
		mockIdx := &mockIndexer{}
		i := NewMockIndexer("", "", nil, nil, mockIdx)

		err := i.Index("path/to/somewhere", nil)
		if err != nil {
			t.Errorf("expected no error, got :%v", err)
		}
	})
	
	t.Run("Index error", func(t *testing.T) {
		mockIdx := &mockIndexer{
			indexError: errors.New("index error"),
		}
		i := NewMockIndexer("", "", nil, nil, mockIdx)

		err := i.Index("path/to/somewhere", nil)
		if err == nil || err.Error() != "index error" {
			t.Errorf("expected 'index error', got: %v", err)
		}
	})
}