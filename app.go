package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"icis/internal/db"
	"icis/internal/iciparser"
	"icis/internal/downloader"
	"icis/internal/extractor"
	"icis/internal/shortcut"
	"icis/internal/uninstaller"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx               context.Context
	database          *db.DB
	uninst            *uninstaller.Uninstaller
	pendingICI        *iciparser.ICIFile
	pendingPath       string
	autoInstall       bool
	pendingLaunchPath string
	pendingLaunchAuto bool
	mu                sync.Mutex
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	database, err := db.Open()
	if err != nil {
		fmt.Println("Failed to open database:", err)
		return
	}
	a.database = database
	a.uninst = uninstaller.New(database)

	filePath, autoInstall := parseArgs(os.Args[1:])
	if filePath != "" {
		a.mu.Lock()
		a.pendingLaunchPath = filePath
		a.pendingLaunchAuto = autoInstall
		a.mu.Unlock()
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.database != nil {
		a.database.Close()
	}
}

func (a *App) LoadICIFile(path string) (*iciparser.ICIFile, error) {
	return a.loadICIFile(path, false)
}

func (a *App) loadICIFile(path string, autoInstall bool) (*iciparser.ICIFile, error) {
	ici, err := iciparser.Parse(path)
	if err != nil {
		return nil, err
	}

	a.mu.Lock()
	a.pendingICI = ici
	a.pendingPath = path
	a.autoInstall = autoInstall
	a.mu.Unlock()

	runtime.EventsEmit(a.ctx, "ici-loaded", ici)

	return ici, nil
}

func (a *App) loadICIFileWithAutoInstall(path string, autoInstall bool) (*iciparser.ICIFile, error) {
	return a.loadICIFile(path, autoInstall)
}

func (a *App) LoadICIContent(content string) (*iciparser.ICIFile, error) {
	ici, err := iciparser.ParseString(content)
	if err != nil {
		return nil, err
	}

	a.mu.Lock()
	a.pendingICI = ici
	a.pendingPath = ""
	a.autoInstall = false
	a.mu.Unlock()

	runtime.EventsEmit(a.ctx, "ici-loaded", ici)

	return ici, nil
}

func (a *App) GetPendingICI() *iciparser.ICIFile {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.pendingICI
}

func (a *App) IsAutoInstall() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.autoInstall
}

type PendingFile struct {
	Path string `json:"path"`
	Auto bool   `json:"auto"`
}

func (a *App) GetPendingFile() PendingFile {
	a.mu.Lock()
	defer a.mu.Unlock()
	pf := PendingFile{Path: a.pendingLaunchPath, Auto: a.pendingLaunchAuto}
	a.pendingLaunchPath = ""
	a.pendingLaunchAuto = false
	return pf
}

func (a *App) InstallApp(iciContent string, installDir string) error {
	ici, err := iciparser.ParseString(iciContent)
	if err != nil {
		return fmt.Errorf("failed to parse .ici: %w", err)
	}

	if installDir != "" {
		ici.InstallDir = "custom"
		ici.CustomDir = installDir
	}

	destPath, err := ici.ResolveInstallPath()
	if err != nil {
		return fmt.Errorf("failed to resolve install path: %w", err)
	}

	runtime.EventsEmit(a.ctx, "install-progress", map[string]string{
		"status":  "downloading",
		"message": "Downloading " + ici.Name + "...",
	})

	tmpDir, err := os.MkdirTemp("", "icis-download-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	downloadedFile, err := downloader.Download(ici.URL, tmpDir, func(downloaded, total int64, percent float64) {
		runtime.EventsEmit(a.ctx, "download-progress", map[string]interface{}{
			"downloaded": downloaded,
			"total":      total,
			"percent":    percent,
		})
	})
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	runtime.EventsEmit(a.ctx, "install-progress", map[string]string{
		"status":  "extracting",
		"message": "Extracting files...",
	})

	extractedFiles, err := extractor.Extract(downloadedFile, destPath, func(file string, current, total int) {
		runtime.EventsEmit(a.ctx, "extract-progress", map[string]interface{}{
			"file":    file,
			"current": current,
			"total":   total,
		})
	})
	if err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	relativeFiles := make([]string, len(extractedFiles))
	for i, f := range extractedFiles {
		rel, err := filepath.Rel(destPath, f)
		if err != nil {
			relativeFiles[i] = filepath.Base(f)
		} else {
			relativeFiles[i] = rel
		}
	}

	runtime.EventsEmit(a.ctx, "install-progress", map[string]string{
		"status":  "shortcut",
		"message": "Creating shortcuts...",
	})

	if ici.Shortcut != "" {
		mainExe := ""
		for _, f := range relativeFiles {
			if strings.HasSuffix(strings.ToLower(f), ".exe") {
				mainExe = f
				break
			}
		}
		if mainExe != "" {
			targetPath := filepath.Join(destPath, mainExe)
			shortcut.CreateDesktopShortcut(ici.Shortcut, targetPath, "")
			shortcut.CreateStartMenuShortcut(ici.Shortcut, targetPath, "")

			if ici.Startup {
				shortcut.CreateStartupShortcut(ici.Shortcut, targetPath, "")
			}
		}
	}

	runtime.EventsEmit(a.ctx, "install-progress", map[string]string{
		"status":  "saving",
		"message": "Saving installation record...",
	})

	var iciSource string
	if a.pendingPath != "" {
		data, _ := os.ReadFile(a.pendingPath)
		iciSource = string(data)
	} else {
		iciSource = iciContent
	}

	app := db.InstalledApp{
		Name:        ici.Name,
		Version:     ici.Version,
		InstallPath: destPath,
		Files:       relativeFiles,
		Shortcut:    ici.Shortcut,
		Startup:     ici.Startup,
		ICISource:   iciSource,
	}

	if err := a.database.SaveApp(app); err != nil {
		return fmt.Errorf("failed to save to database: %w", err)
	}

	runtime.EventsEmit(a.ctx, "install-complete", map[string]string{
		"name":  ici.Name,
		"path":  destPath,
		"message": fmt.Sprintf("%s installed successfully!", ici.Name),
	})

	return nil
}

func (a *App) UninstallApp(appName string) error {
	runtime.EventsEmit(a.ctx, "uninstall-progress", map[string]string{
		"status":  "removing",
		"message": "Uninstalling " + appName + "...",
	})

	if err := a.uninst.Uninstall(appName); err != nil {
		return fmt.Errorf("uninstall failed: %w", err)
	}

	runtime.EventsEmit(a.ctx, "uninstall-complete", map[string]string{
		"name":    appName,
		"message": appName + " has been uninstalled.",
	})

	return nil
}

func (a *App) ListInstalledApps() ([]db.InstalledApp, error) {
	if a.database == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return a.database.ListApps()
}

func (a *App) GetApp(name string) (*db.InstalledApp, error) {
	if a.database == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return a.database.GetApp(name)
}

func (a *App) OpenFile() (string, error) {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select .ici File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "ICIS Files (*.ici)",
				Pattern:     "*.ici",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

func (a *App) SelectDirectory() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Install Directory",
	})
	if err != nil {
		return "", err
	}
	return dir, nil
}
