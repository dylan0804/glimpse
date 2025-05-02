package main

import (
	"context"
	"glimpse/screenshots"
)

// App struct
type App struct {
	ctx context.Context
	screenshotService screenshots.Service
	ocr screenshots.OCRProvider
}

// NewApp creates a new App application struct
func NewApp() *App {
	d := screenshots.NewDirProvider()
	o := screenshots.NewOCRProvider()

	return &App{
		screenshotService: screenshots.NewScreenshotService(d, o),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// func (a *App) shutdown(ctx context.Context) {
// 	a.ocr.Close()
// }

func (a *App) ScanScreenshots() error {
	a.screenshotService.ScanAndIndex()
	return nil
}


