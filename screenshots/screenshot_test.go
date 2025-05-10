package screenshots

// import (
// 	"io/fs"
// 	"testing"
// )

// type fakeDir struct {
// 	homeDir string
// 	entries []fs.DirEntry
// 	errHomeDir error
// 	errReadDir error
// }
// func (f *fakeDir) GetHomeDir() (string, error) { return f.homeDir, f.errHomeDir}
// func (f *fakeDir) ReadDir(string) ([]fs.DirEntry, error) { return f.entries, f.errReadDir}

// type fakeOcr struct {}
// func (f *fakeOcr) ExtractText(path string) (string, error) { return "extracted text", nil}
// func (f *fakeOcr) Close()

// type fakeEntry struct { name string }
// func (f fakeEntry) Name() string               { return f.name }
// func (fakeEntry) IsDir() bool                  { return false }
// func (fakeEntry) Type() fs.FileMode            { return 0 }
// func (fakeEntry) Info() (fs.FileInfo, error)   { return nil, nil }

// func TestScanAndIndex(t *testing.T) {
// 	dir := &fakeDir{
//         homeDir: "/home",
// 		entries: []fs.DirEntry{
// 			fakeEntry{name: "img1.png"},
// 			fakeEntry{name: "img2.png"},
// 		},
//     }
// 	ocr := &fakeOcr{}

// 	svc := NewScreenshotService(dir, ocr)
// 	if 	err := svc.ScanAndIndex(); err != nil {
// 		t.Fatal("expected no error; got", err)
// 	}

// }