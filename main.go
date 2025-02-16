package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "go-organizer",
		Width:  420,
		Height: 460,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 48, G: 48, B: 48, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		DisableResize:   false,
		OnBeforeClose:   app.beforeClose,
		Frameless:       true,
		CSSDragProperty: "widows",
		CSSDragValue:    "1",
		MinWidth:        1,
		MinHeight:       1,
		MaxWidth:        420,
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
