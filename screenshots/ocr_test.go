package screenshots

import (
	"errors"
	"os"
	"testing"
)

type mockCmdRunner struct {
	output []byte
	cmdError error
}

// mockFile implements the File interface
type mockFile struct {
	name string
}

type mockFileSystem struct {
	createTempError error
	writeFileError error
	chmodError error
	tempFile File
}

func (m *mockFileSystem) CreateTemp(dir string, pattern string) (File, error) {
	if m.createTempError != nil {
		return nil, m.createTempError
	}

	return m.tempFile, nil
}

func (m *mockFileSystem) WriteFile(file File, b []byte) (int, error) {
	if m.writeFileError != nil {
		return 0, m.writeFileError
	}
	
	return len(b), nil
}

func (m *mockFileSystem) Chmod(name string, mode os.FileMode) error {
	if m.chmodError != nil {
		return m.chmodError
	}
	return nil
}

func NewMockOCR(ocrBinary []byte, ocrBinaryPath string, fs fileSystem, cmdRunner commandRunner) *OCR {
	return &OCR{
		ocrBinary: ocrBinary,
		ocrBinaryPath: ocrBinaryPath,
		fs: fs,
		cmdRunner: cmdRunner,
	}
}

func (m *mockFile) Name() string {
	return m.name
}

func (m *mockFile) Close() error {
	return nil
}

func (m *mockFile) Write(b []byte) (int, error) {
	return len(b), nil
}

func TestWriteOCRHelper(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockFile := &mockFile{name: "/tmp/test-ocr-binary"}
		mockFileSystem := &mockFileSystem{tempFile: mockFile}

		ocr := NewMockOCR([]byte("binary content"), "", mockFileSystem, nil)
		
		path, err := ocr.WriteOCRHelper()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if path != "/tmp/test-ocr-binary" {
			t.Errorf("expected path to be '/tmp/test-ocr-binary', got: %s", path)
		}
	})

	t.Run("CreateTemp error", func(t *testing.T) {
		err := errors.New("create temp error")
		mockFileSystem := &mockFileSystem{
			createTempError: err,
		}
		ocr := NewMockOCR(nil, "", mockFileSystem, nil)

		_, err = ocr.WriteOCRHelper()
		if err == nil || err.Error() != "create temp error" {
			t.Errorf("expected error 'create temp error', got: %v", err)
		}
	})

	t.Run("WriteFile error", func(t *testing.T) {
		err := errors.New("write file error")
		mockFile := &mockFile{name: "/tmp/test-ocr-binary"}
		mockFileSystem := &mockFileSystem{
			writeFileError: err,
			tempFile: mockFile,
		}
		ocr := NewMockOCR([]byte("binary content"), "", mockFileSystem, nil)

		_, err = ocr.WriteOCRHelper()
		if err == nil || err.Error() != "write file error" {
			t.Errorf("expected error 'write file error', got: %v", err)
		}
	})

	t.Run("Chmod error", func(t *testing.T) {
		err := errors.New("Chmod error")
		mockFile := &mockFile{name: "/tmp/test-ocr-binary"}
		mockFileSystem := &mockFileSystem{
			writeFileError: err,
			tempFile: mockFile,
		}
		ocr := NewMockOCR(nil, "", mockFileSystem, nil)

		_, err = ocr.WriteOCRHelper()
		if err == nil || err.Error() != "Chmod error" {
			t.Errorf("expected error 'Chmod error', got: %v", err)
		}
	})
}

func (m *mockCmdRunner) Command(name string, arg ...string) ([]byte, error) {
	if m.cmdError != nil {
		return nil, m.cmdError
	}

	return m.output, nil
}

func TestExtractText (t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		cmdRunner := &mockCmdRunner{
			output: []byte("text extracted from ss"),
		}
		ocr := NewMockOCR(nil, "", nil, cmdRunner)

		ocr.ocrBinaryPath = "/path/to/binary"

		extractedText, err := ocr.ExtractText("/path/to/screenshot")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if extractedText != "text extracted from ss" {
			t.Errorf("expected 'text extracted from ss', got: %s", extractedText)
		}
	})

	t.Run("Cmd error", func(t *testing.T) {
		cmdRunner := &mockCmdRunner{
			cmdError: errors.New("cmd error"),
		}
		ocr := NewMockOCR(nil, "", nil, cmdRunner)

		ocr.ocrBinaryPath = "/path/to/binary"

		_, err := ocr.ExtractText("/path/to/screenshot")
		if err == nil || err.Error() != "cmd error" {
			t.Errorf("expected 'cmd error', got: %v", err)
		}
	})
}


