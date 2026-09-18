package shortcut

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func getStartupFolder() (string, error) {
	appData, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appData, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup"), nil
}

func getStartMenuProgramsFolder() (string, error) {
	appData, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appData, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs"), nil
}

func CreateDesktopShortcut(name, targetPath, iconPath string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("shortcuts only supported on Windows")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	desktopPath := filepath.Join(home, "Desktop")
	shortcutPath := filepath.Join(desktopPath, name+".lnk")

	return createShortcutCOM(shortcutPath, targetPath, filepath.Dir(targetPath), iconPath)
}

func CreateStartMenuShortcut(name, targetPath, iconPath string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("shortcuts only supported on Windows")
	}

	programsDir, err := getStartMenuProgramsFolder()
	if err != nil {
		return err
	}

	shortcutDir := filepath.Join(programsDir, "ICIS")
	if err := os.MkdirAll(shortcutDir, 0o755); err != nil {
		return err
	}

	shortcutPath := filepath.Join(shortcutDir, name+".lnk")
	return createShortcutCOM(shortcutPath, targetPath, filepath.Dir(targetPath), iconPath)
}

func CreateStartupShortcut(name, targetPath, iconPath string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("shortcuts only supported on Windows")
	}

	startupDir, err := getStartupFolder()
	if err != nil {
		return err
	}

	shortcutPath := filepath.Join(startupDir, name+".lnk")
	return createShortcutCOM(shortcutPath, targetPath, filepath.Dir(targetPath), iconPath)
}

func RemoveDesktopShortcut(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, "Desktop", name+".lnk")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func RemoveStartMenuShortcut(name string) error {
	programsDir, err := getStartMenuProgramsFolder()
	if err != nil {
		return err
	}
	path := filepath.Join(programsDir, "ICIS", name+".lnk")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	dir := filepath.Join(programsDir, "ICIS")
	entries, _ := os.ReadDir(dir)
	if len(entries) == 0 {
		os.Remove(dir)
	}

	return nil
}

func RemoveStartupShortcut(name string) error {
	startupDir, err := getStartupFolder()
	if err != nil {
		return err
	}
	path := filepath.Join(startupDir, name+".lnk")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func createShortcutCOM(shortcutPath, targetPath, workingDir, iconPath string) error {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_DISABLE_OLE1DDE)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return fmt.Errorf("failed to create WScript.Shell: %w", err)
	}
	defer unknown.Release()

	shell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("failed to query IDispatch: %w", err)
	}
	defer shell.Release()

	result, err := oleutil.CallMethod(shell, "CreateShortcut", shortcutPath)
	if err != nil {
		return fmt.Errorf("failed to create shortcut: %w", err)
	}

	link := result.ToIDispatch()
	defer link.Release()

	oleutil.PutProperty(link, "TargetPath", targetPath)
	oleutil.PutProperty(link, "WorkingDirectory", workingDir)
	if iconPath != "" {
		oleutil.PutProperty(link, "IconLocation", iconPath)
	}

	_, err = oleutil.CallMethod(link, "Save")
	if err != nil {
		return fmt.Errorf("failed to save shortcut: %w", err)
	}

	return nil
}
