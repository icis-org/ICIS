package pack

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var magic = []byte("ICIS")

// Pack creates a standalone installer by concatenating: [stub][ici][installer][trailer]
// No Go compilation or network needed at pack time.
func Pack(iciPath string, outputPath string, stubBinary []byte, installerBinary []byte) error {
	// Read .ici content
	iciBytes, err := os.ReadFile(iciPath)
	if err != nil {
		return fmt.Errorf("failed to read .ici file: %w", err)
	}

	// Ensure .exe extension
	if runtime.GOOS == "windows" && !strings.HasSuffix(outputPath, ".exe") {
		outputPath += ".exe"
	}

	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		absOutput = outputPath
	}

	// Build the output file: [stub][ici][installer][iciLen:u32][installerLen:u32]["ICIS"]
	f, err := os.Create(absOutput)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	// 1. Write stub binary
	if _, err := f.Write(stubBinary); err != nil {
		return fmt.Errorf("failed to write stub: %w", err)
	}

	// 2. Write .ici content
	if _, err := f.Write(iciBytes); err != nil {
		return fmt.Errorf("failed to write .ici: %w", err)
	}

	// 3. Write ICIS installer
	if _, err := f.Write(installerBinary); err != nil {
		return fmt.Errorf("failed to write installer: %w", err)
	}

	// 4. Write trailer: [iciLen uint32][installerLen uint32][magic "ICIS"]
	trailer := make([]byte, 12)
	binary.LittleEndian.PutUint32(trailer[0:4], uint32(len(iciBytes)))
	binary.LittleEndian.PutUint32(trailer[4:8], uint32(len(installerBinary)))
	copy(trailer[8:12], magic)

	if _, err := f.Write(trailer); err != nil {
		return fmt.Errorf("failed to write trailer: %w", err)
	}

	f.Close()

	// Print result
	info, _ := os.Stat(absOutput)
	sizeMB := float64(info.Size()) / 1024 / 1024
	fmt.Printf("Standalone installer created: %s (%.1f MB)\n", absOutput, sizeMB)
	return nil
}

// FindInstaller looks for the ICIS NSIS installer relative to the ICIS executable.
func FindInstaller() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)

	// Look for installer in same directory as icis.exe
	patterns := []string{
		"icis-amd64-installer.exe",
		"icis-*-installer.exe",
		"installer.exe",
	}

	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(exeDir, pattern))
		if len(matches) > 0 {
			return matches[0], nil
		}
	}

	// Also check build/bin relative to source (for dev)
	cwd, _ := os.Getwd()
	buildPaths := []string{
		filepath.Join(cwd, "build", "bin"),
		filepath.Join(exeDir, "..", "build", "bin"),
	}

	for _, dir := range buildPaths {
		for _, pattern := range patterns {
			matches, _ := filepath.Glob(filepath.Join(dir, pattern))
			if len(matches) > 0 {
				return matches[0], nil
			}
		}
	}

	return "", fmt.Errorf("ICIS installer not found — expected icis-amd64-installer.exe next to icis.exe")
}
