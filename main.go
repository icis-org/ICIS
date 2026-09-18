package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "ICIS - Isam's Configurable Install System",
		Width:     800,
		Height:    600,
		MinWidth:  700,
		MinHeight: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 18, G: 18, B: 24, A: 1},
		OnStartup:        app.startup,
		OnBeforeClose: func(ctx context.Context) bool {
			return false
		},
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "e3984e08-28dc-4e3d-b70a-45e961589cdc",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				if len(secondInstanceData.Args) > 0 {
					for _, arg := range secondInstanceData.Args {
						if len(arg) > 4 && arg[len(arg)-4:] == ".ici" {
							app.loadICIFile(arg)
						}
					}
				}
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
