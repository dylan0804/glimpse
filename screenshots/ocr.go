package screenshots

import (
	_ "embed"
	"os"
	"os/exec"
	"strings"
)

type OCRProvider interface {
	ExtractText(path string) (string, error)
	WriteOCRHelper() error
}

type OCR struct {
	ocrBinary []byte
	ocrBinaryPath string
}

func NewOCRProvider(ocrBinary []byte) OCRProvider {
	return &OCR{
		ocrBinary: ocrBinary,
	}
}

func (o *OCR) ExtractText(path string) (string, error) {
	out, err := exec.Command(o.ocrBinaryPath, path).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func (o *OCR) WriteOCRHelper() error {
	tmpFile, err := os.CreateTemp("", "ocr-helper-*")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(o.ocrBinary); err != nil {
		return err
	}

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return err
	}

	o.ocrBinaryPath = tmpFile.Name()

	return nil
}