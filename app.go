package main

import (
	_ "embed"

	"context"
	"glimpse/screenshots"
)

//go:embed ocr-helper
var ocrHelper []byte

// App struct
type App struct {
	ctx context.Context
	screenshotService screenshots.Service
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	d := screenshots.NewDirProvider()
	o := screenshots.NewOCRProvider(ocrHelper)
	i := screenshots.NewIndexer()

	a.screenshotService = screenshots.NewScreenshotService(d, o, i, a.ctx)
}

func (a *App) shutdown(ctx context.Context) {
	a.screenshotService.Shutdown()
}

func (a *App) ScanScreenshots() error {
	err := a.screenshotService.ScanAndIndex()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) SearchScreenshots(query string) error {
	err := a.screenshotService.Search(query)
	if err != nil {
		return err
	}
	
	return nil
}




