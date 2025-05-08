package screenshots

import (
	_ "embed"
	"os"
	"os/exec"
	"strings"
)

type OCRProvider interface {
	ExtractText(path string) (string, error)
	WriteOCRHelper() (string, error)
}

type File interface {
	Name() string
	Write(b []byte) (int, error)
	Close() error
}

type fileSystem interface {
	CreateTemp(dir string, pattern string) (File, error)
	WriteFile(file File, b []byte) (int, error)
	Chmod(name string, mode os.FileMode) error
}

type commandRunner interface {
	Command(name string, arg ...string) ([]byte, error)
}

type OCR struct {
	ocrBinary []byte
	ocrBinaryPath string
	fs fileSystem
	cmdRunner commandRunner
}

type realFileSystem struct {}

func (r *realFileSystem) CreateTemp(dir string, pattern string) (File, error) {
	return os.CreateTemp(dir, pattern)
}

func (r *realFileSystem) WriteFile(file File, b []byte) (int, error) {
	return file.Write(b)
}

func (r *realFileSystem) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

type realCommandRunner struct {}

func (c *realCommandRunner) Command(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).Output()
}

func NewOCRProvider(ocrBinary []byte) *OCR {
	return &OCR{
		ocrBinary: ocrBinary,
		fs: &realFileSystem{},
		cmdRunner: &realCommandRunner{},
	}
}

func (o *OCR) ExtractText(path string) (string, error) {
	out, err := o.cmdRunner.Command(o.ocrBinaryPath, path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func (o *OCR) WriteOCRHelper() (string, error) {
	tmpFile, err := o.fs.CreateTemp("", "ocr-helper-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := o.fs.WriteFile(tmpFile, o.ocrBinary); err != nil {
		return "", err
	}

	if err := o.fs.Chmod(tmpFile.Name(), 0755); err != nil {
		return "", err
	}

	o.ocrBinaryPath = tmpFile.Name()

	return o.ocrBinaryPath, nil
}
