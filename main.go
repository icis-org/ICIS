package main

import (
	"context"
	"embed"
	"strings"

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
					filePath, autoInstall := parseArgs(secondInstanceData.Args)
					if filePath != "" {
						app.loadICIFileWithAutoInstall(filePath, autoInstall)
					}
				}
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func parseArgs(args []string) (filePath string, autoInstall bool) {
	for i, arg := range args {
		if arg == "--install" && i+1 < len(args) {
			next := args[i+1]
			if strings.HasSuffix(strings.ToLower(next), ".ici") {
				return next, true
			}
		}
	}

	for _, arg := range args {
		if strings.HasSuffix(strings.ToLower(arg), ".ici") {
			return arg, false
		}
	}

	return "", false
}
