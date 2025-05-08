package screenshots

import (
	"errors"
	"io/fs"
	"testing"
)

type mockDir struct {
	mockHomeDir string
	mockEntries []fs.DirEntry
	mockError error
}

func newMockDirProvider(homeDir string, entries []fs.DirEntry, mockError error) DirProvider {
	return &mockDir{
		mockHomeDir: homeDir,
		mockEntries: entries,
		mockError: mockError,
	}
}

func (m *mockDir) GetHomeDir() (string, error) {
	if m.mockError != nil {
		return "", m.mockError
	}
	
	return m.mockHomeDir, nil
}

func (m *mockDir) ReadDir(path string) ([]fs.DirEntry, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}

	return m.mockEntries, nil
}

type MockDirEntry struct {
	name     string
	isDir    bool
	fileMode fs.FileMode
	fileInfo fs.FileInfo
	err      error
}

func (m *MockDirEntry) Name() string {
	return m.name
}

func (m *MockDirEntry) IsDir() bool {
	return m.isDir
}

func (m *MockDirEntry) Type() fs.FileMode {
	return m.fileMode
}

func (m *MockDirEntry) Info() (fs.FileInfo, error) {
	return m.fileInfo, m.err
}

func TestMockDirProvider(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockEntries := []fs.DirEntry{
			&MockDirEntry{name: "file1.txt", isDir: false},
			&MockDirEntry{name: "file2.txt", isDir: false},
			&MockDirEntry{name: "dir1", isDir: true},
		}

		mockProvider := newMockDirProvider("/mock/home", mockEntries, nil)

		// test GetHomeDir
		homeDir, err := mockProvider.GetHomeDir()
		if err != nil {
			t.Errorf("expected no error, got :%v", err)
		}
		if homeDir != "/mock/home" {
			t.Errorf("expected /mock/home, got %s", homeDir)
		}

		// test ReadDir
		entries, err := mockProvider.ReadDir("/some/path")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if entries[0].Name() != "file1.txt" {
			t.Errorf("expected first entry name to be 'file1.txt', got: %s", entries[0].Name())
		}
		if entries[2].IsDir() != true {
			t.Errorf("expected third entry to be a directory")
		}
	})

	t.Run("ErrorGettingHomeDir", func(t *testing.T) {
		mockError := errors.New("home dir error")
		mockProvider := newMockDirProvider("", nil, mockError)

		_, err := mockProvider.GetHomeDir()
		if err == nil || err.Error() != "home dir error" {
			t.Errorf("expected 'home dir error', got %v", err)
		}
	})

	t.Run("ErrorReadingDir", func(t *testing.T) {
		mockError := errors.New("error reading dir")
		mockProvider := newMockDirProvider("", nil, mockError)

		_, err := mockProvider.ReadDir("/ts/pmo")
		if err == nil || err.Error() != "error reading dir" {
			t.Errorf("expected 'error reading dir', got: %v", err)
		}
	})
}