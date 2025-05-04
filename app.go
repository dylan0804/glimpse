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
	cancel context.CancelFunc
	screenshotService screenshots.Service
}

// NewApp creates a new App application struct
func NewApp() *App {
	d := screenshots.NewDirProvider()
	o := screenshots.NewOCRProvider(ocrHelper)
	i := screenshots.NewIndexer()

	return &App{
		screenshotService: screenshots.NewScreenshotService(d, o, i),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context, cancel context.CancelFunc) {
	a.ctx = ctx
	a.cancel = cancel
}

func (a *App) shutdown(_ context.Context) {
	a.cancel()
}

func (a *App) ScanScreenshots() error {
	err := a.screenshotService.ScanAndIndex(a.ctx)
	if err != nil {
		return err
	}

	return nil
}


