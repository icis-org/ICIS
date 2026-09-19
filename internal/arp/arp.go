package arp

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"icis/internal/db"

	"golang.org/x/sys/windows/registry"
)

const uninstallKeyPath = `Software\Microsoft\Windows\CurrentVersion\Uninstall`
const keyPrefix = "ICIS_"

var sanitizeRe = regexp.MustCompile(`[^A-Za-z0-9_\-]`)

func sanitizeKeyName(name string) string {
	return keyPrefix + sanitizeRe.ReplaceAllString(name, "_")
}

func Register(app db.InstalledApp) {
	keyName := sanitizeKeyName(app.Name)

	k, _, err := registry.CreateKey(registry.CURRENT_USER, uninstallKeyPath+`\`+keyName, registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		return
	}
	defer k.Close()

	k.SetStringValue("DisplayName", app.Name)
	if app.Version != "" {
		k.SetStringValue("DisplayVersion", app.Version)
	}
	k.SetStringValue("InstallLocation", app.InstallPath)
	k.SetStringValue("Publisher", "ICIS")

	icisExe := icisExePath()
	uninstallCmd := fmt.Sprintf(`"%s" --uninstall "%s"`, icisExe, app.Name)
	k.SetStringValue("UninstallString", uninstallCmd)
	k.SetStringValue("QuietUninstallString", uninstallCmd+" /S")
	k.SetStringValue("DisplayIcon", displayIcon(app))

	k.SetDWordValue("NoModify", 1)
	k.SetDWordValue("NoRepair", 1)

	if sizeKB := estimateSizeKB(app.InstallPath); sizeKB > 0 {
		k.SetDWordValue("EstimatedSize", uint32(sizeKB))
	}
}

func Unregister(name string) {
	keyName := sanitizeKeyName(name)
	registry.DeleteKey(registry.CURRENT_USER, uninstallKeyPath+`\`+keyName)
}

func Backfill(database *db.DB) {
	apps, err := database.ListApps()
	if err != nil {
		return
	}

	registered := make(map[string]bool)

	k, err := registry.OpenKey(registry.CURRENT_USER, uninstallKeyPath, registry.READ)
	if err != nil {
		for _, app := range apps {
			Register(app)
		}
		return
	}
	defer k.Close()

	subkeys, _ := k.ReadSubKeyNames(-1)
	for _, sub := range subkeys {
		if strings.HasPrefix(sub, keyPrefix) {
			registered[sub] = true
		}
	}

	for _, app := range apps {
		keyName := sanitizeKeyName(app.Name)
		if !registered[keyName] {
			Register(app)
		} else {
			delete(registered, keyName)
		}
	}

	for orphanKey := range registered {
		registry.DeleteKey(registry.CURRENT_USER, uninstallKeyPath+`\`+orphanKey)
	}
}

func SweepOrphans() {
	k, err := registry.OpenKey(registry.CURRENT_USER, uninstallKeyPath, registry.READ)
	if err != nil {
		return
	}
	defer k.Close()

	subkeys, _ := k.ReadSubKeyNames(-1)
	for _, sub := range subkeys {
		if strings.HasPrefix(sub, keyPrefix) {
			registry.DeleteKey(registry.CURRENT_USER, uninstallKeyPath+`\`+sub)
		}
	}
}

func displayIcon(app db.InstalledApp) string {
	for _, sl := range app.Shortcuts {
		full := filepath.Join(app.InstallPath, sl.Exe)
		if _, err := os.Stat(full); err == nil {
			return full
		}
	}
	return icisExePath()
}

func icisExePath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return exe
}

func estimateSizeKB(dir string) int64 {
	var total int64
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total / 1024
}
