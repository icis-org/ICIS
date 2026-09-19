package main

import (
	"context"
	"embed"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"icis/internal/arp"
	"icis/internal/db"
	"icis/internal/pack"
	"icis/internal/uninstaller"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed internal/pack/stub.exe
var stubBinary []byte

var consoleHandle uintptr

func attachConsole() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	allocConsole := kernel32.NewProc("AllocConsole")
	allocConsole.Call()
	getStdHandle := kernel32.NewProc("GetStdHandle")
	hOut, _, _ := getStdHandle.Call(^uintptr(11) + 1)
	if hOut != 0 && hOut != uintptr(0xFFFFFFFFFFFFFFFF) {
		consoleHandle = hOut
	}
}

func consolePrint(msg string) {
	if consoleHandle != 0 {
		b := []byte(msg)
		syscall.Write(syscall.Handle(consoleHandle), b)
	}
}

type consoleWriter struct{}

func (consoleWriter) Write(p []byte) (int, error) {
	if consoleHandle != 0 {
		syscall.Write(syscall.Handle(consoleHandle), p)
	}
	return len(p), nil
}

func main() {
	if headlessPack() {
		return
	}
	if headlessUninstall() {
		return
	}

	app := NewApp()

	err := wails.Run(&options.App{
		Title:      "ICIS - Isam's Configurable Install System",
		Width:      960,
		Height:     640,
		MinWidth:   800,
		MinHeight:  540,
		Frameless:  true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		OnBeforeClose: func(ctx context.Context) bool {
			return app.shouldBlockClose()
		},
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Acrylic,
			Theme:                windows.Dark,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "e3984e08-28dc-4e3d-b70a-45e961589cdc",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				if len(secondInstanceData.Args) > 0 {
					filePath, autoInstall, protocolURL := parseArgs(secondInstanceData.Args)
					if protocolURL != "" {
						app.handleProtocolURL(protocolURL)
					} else if filePath != "" {
						app.loadICIFileWithAutoInstall(filePath, autoInstall)
					}
				}
				app.focusWindow()
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func headlessUninstall() bool {
	for i, arg := range os.Args[1:] {
		if arg == "--uninstall" && i+1 < len(os.Args)-1 {
			attachConsole()
			appName := os.Args[i+2]
			database, err := db.Open()
			if err != nil {
				consolePrint("Error: cannot open database: " + err.Error() + "\n")
				os.Exit(1)
			}
			inst := uninstaller.New(database)
			inst.Uninstall(appName)
			database.Close()
			arp.SweepOrphans()
			return true
		}
	}
	return false
}

func headlessPack() bool {
	for i, arg := range os.Args[1:] {
		if arg == "pack" {
			attachConsole()
			pack.SetOutput(consoleWriter{})
			args := os.Args[i+2:]
			if len(args) < 1 {
				consolePrint("Usage: icis pack <input.ici> [-o output.exe]\n")
				os.Exit(1)
			}
			iciPath := args[0]
			outputPath := strings.TrimSuffix(iciPath, filepath.Ext(iciPath)) + "-installer.exe"
			for j, a := range args {
				if a == "-o" && j+1 < len(args) {
					outputPath = args[j+1]
				}
			}

			installerPath, err := pack.FindInstaller()
			if err != nil {
				consolePrint("Error: " + err.Error() + "\n")
				os.Exit(1)
			}
			installerBinary, err := os.ReadFile(installerPath)
			if err != nil {
				consolePrint("Error: cannot read ICIS installer: " + err.Error() + "\n")
				os.Exit(1)
			}

			if err := pack.Pack(iciPath, outputPath, stubBinary, installerBinary); err != nil {
				consolePrint("Error: " + err.Error() + "\n")
				os.Exit(1)
			}
			return true
		}
	}
	return false
}

func parseArgs(args []string) (filePath string, autoInstall bool, protocolURL string) {
	for _, arg := range args {
		if strings.HasPrefix(strings.ToLower(arg), "icis://") {
			parsed, err := url.Parse(arg)
			if err != nil {
				continue
			}
			if strings.ToLower(parsed.Host) != "install" {
				continue
			}
			iciURL := parsed.Query().Get("ici")
			if iciURL == "" {
				continue
			}
			if !strings.HasPrefix(strings.ToLower(iciURL), "https://") {
				continue
			}
			return "", false, iciURL
		}
	}

	for i, arg := range args {
		if arg == "--install" && i+1 < len(args) {
			next := args[i+1]
			if strings.HasSuffix(strings.ToLower(next), ".ici") {
				return next, true, ""
			}
		}
	}

	for _, arg := range args {
		if strings.HasSuffix(strings.ToLower(arg), ".ici") {
			return arg, false, ""
		}
	}

	return "", false, ""
}
