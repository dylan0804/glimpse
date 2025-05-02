package screenshots

import (
	"os/exec"
	"strings"
)

type OCRProvider interface {
	ExtractText(path string) (string, error)
}

type OCR struct {}

func NewOCRProvider() OCRProvider {
	return &OCR{}
}

func (o *OCR) ExtractText(path string) (string, error) {
	out, err := exec.Command("../ocr-helper/ocr-helper", path).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}