package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var magic = []byte("ICIS")

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error: cannot determine executable path:", err)
		waitForExit()
		return
	}

	f, err := os.Open(exePath)
	if err != nil {
		fmt.Println("Error: cannot open executable:", err)
		waitForExit()
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		fmt.Println("Error: cannot stat executable:", err)
		waitForExit()
		return
	}

	fileSize := stat.Size()
	if fileSize < 12 {
		fmt.Println("Error: invalid installer package (too small)")
		waitForExit()
		return
	}

	// Read trailer: [iciLen uint32][installerLen uint32][magic 4 bytes]
	trailer := make([]byte, 12)
	_, err = f.ReadAt(trailer, fileSize-12)
	if err != nil {
		fmt.Println("Error: cannot read installer trailer:", err)
		waitForExit()
		return
	}

	if string(trailer[8:12]) != string(magic) {
		fmt.Println("Error: invalid installer package (bad magic)")
		waitForExit()
		return
	}

	iciLen := binary.LittleEndian.Uint32(trailer[0:4])
	installerLen := binary.LittleEndian.Uint32(trailer[4:8])

	dataStart := fileSize - 12 - int64(installerLen) - int64(iciLen)
	if dataStart < 0 {
		fmt.Println("Error: invalid installer package (corrupted sizes)")
		waitForExit()
		return
	}

	// Read .ici content
	iciBytes := make([]byte, iciLen)
	_, err = f.ReadAt(iciBytes, dataStart)
	if err != nil {
		fmt.Println("Error: cannot read embedded config:", err)
		waitForExit()
		return
	}

	// Read ICIS installer
	installerBytes := make([]byte, installerLen)
	_, err = f.ReadAt(installerBytes, dataStart+int64(iciLen))
	if err != nil {
		fmt.Println("Error: cannot read embedded installer:", err)
		waitForExit()
		return
	}

	f.Close()

	// Parse app name from .ici content
	appName := "Application"
	for _, line := range strings.Split(string(iciBytes), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "name:") {
			appName = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			break
		}
	}

	fmt.Printf("=== %s Installer ===\n", appName)
	fmt.Println()

	// Find or install ICIS
	icisPath := findICIS()
	if icisPath == "" {
		fmt.Println("ICIS is not installed. Installing ICIS first...")
		if err := installICIS(installerBytes); err != nil {
			fmt.Printf("Failed to install ICIS: %v\n", err)
			fmt.Println("Please install ICIS manually from https://github.com/icis-org/icis/releases")
			waitForExit()
			return
		}
		icisPath = findICIS()
		if icisPath == "" {
			fmt.Println("ICIS installation failed. Please install manually.")
			waitForExit()
			return
		}
		fmt.Println("ICIS installed successfully.")
		fmt.Println()
	}

	// Write .ici to temp
	tmpDir := os.TempDir()
	iciFile := filepath.Join(tmpDir, "icis-pack.ici")
	if err := os.WriteFile(iciFile, iciBytes, 0644); err != nil {
		fmt.Printf("Failed to write config: %v\n", err)
		waitForExit()
		return
	}
	defer os.Remove(iciFile)

	// Run ICIS with --install
	fmt.Printf("Launching ICIS installer for %s...\n", appName)
	cmd := exec.Command(icisPath, "--install", iciFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Installation failed: %v\n", err)
	}

	waitForExit()
}

func findICIS() string {
	candidates := []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "ICIS", "icis.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "ICIS", "icis.exe"),
	}

	if p, err := exec.LookPath("icis.exe"); err == nil {
		candidates = append([]string{p}, candidates...)
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func installICIS(installerBytes []byte) error {
	tmpDir := os.TempDir()
	installerPath := filepath.Join(tmpDir, "icis-installer.exe")

	if err := os.WriteFile(installerPath, installerBytes, 0755); err != nil {
		return fmt.Errorf("failed to write installer: %w", err)
	}
	defer os.Remove(installerPath)

	cmd := exec.Command(installerPath, "/S", "/CURRENTUSER")
	if err := cmd.Run(); err != nil {
		cmd2 := exec.Command(installerPath, "/CURRENTUSER")
		cmd2.Run()
	}

	return nil
}

func waitForExit() {
	fmt.Println()
	fmt.Println("Press Enter to exit...")
	var b [1]byte
	os.Stdin.Read(b[:])
}
