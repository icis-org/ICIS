package uninstaller

import (
	"fmt"
	"os"
	"path/filepath"

	"icis/internal/db"
	"icis/internal/shortcut"
)

type Uninstaller struct {
	database *db.DB
}

func New(database *db.DB) *Uninstaller {
	return &Uninstaller{database: database}
}

func (u *Uninstaller) Uninstall(appName string) error {
	app, err := u.database.GetApp(appName)
	if err != nil {
		return fmt.Errorf("app not found: %w", err)
	}

	if err := u.removeFiles(app); err != nil {
		return fmt.Errorf("failed to remove files: %w", err)
	}

	if len(app.Shortcuts) > 0 {
		for i, sl := range app.Shortcuts {
			shortcut.RemoveDesktopShortcut(sl.Name)
			shortcut.RemoveStartMenuShortcut(sl.Name)
			if app.Startup && i == 0 {
				shortcut.RemoveStartupShortcut(sl.Name)
			}
		}
	} else if app.Shortcut != "" {
		shortcut.RemoveDesktopShortcut(app.Shortcut)
		shortcut.RemoveStartMenuShortcut(app.Shortcut)
		if app.Startup {
			shortcut.RemoveStartupShortcut(app.Shortcut)
		}
	}

	if err := u.database.DeleteApp(appName); err != nil {
		return fmt.Errorf("failed to remove from database: %w", err)
	}

	return nil
}

func (u *Uninstaller) removeFiles(app *db.InstalledApp) error {
	if _, err := os.Stat(app.InstallPath); os.IsNotExist(err) {
		return nil
	}

	if len(app.Files) > 0 {
		for _, file := range app.Files {
			fullPath := filepath.Join(app.InstallPath, file)
			if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}

	entries, err := os.ReadDir(app.InstallPath)
	if err != nil {
		return nil
	}

	if len(entries) == 0 {
		if err := os.Remove(app.InstallPath); err != nil {
			return err
		}
	}

	return nil
}

func (u *Uninstaller) GetInstalledApps() ([]db.InstalledApp, error) {
	return u.database.ListApps()
}
