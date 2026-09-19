package pack

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var output io.Writer = os.Stderr

func SetOutput(w io.Writer) {
	output = w
}

func printf(format string, args ...interface{}) {
	fmt.Fprintf(output, format, args...)
}

func println(args ...interface{}) {
	fmt.Fprintln(output, args...)
}

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
	printf("Standalone installer created: %s (%.1f MB)\n", absOutput, sizeMB)
	return nil
}

// FindInstaller looks for the ICIS NSIS installer relative to the ICIS executable.
// If not found locally, downloads from GitHub releases.
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

	// Not found locally — download from GitHub releases
	return downloadInstaller(exeDir)
}

const icisReleaseURL = "https://github.com/icis-org/icis/releases/latest/download/icis-amd64-installer.exe"

func downloadInstaller(targetDir string) (string, error) {
	destPath := filepath.Join(targetDir, "icis-amd64-installer.exe")

	println("ICIS installer not found locally. Downloading from GitHub...")

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false}}
	client := &http.Client{Transport: tr, Timeout: 120 * time.Second}

	resp, err := client.Get(icisReleaseURL)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	total := resp.ContentLength
	written := int64(0)
	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			f.Write(buf[:n])
			written += int64(n)
			if total > 0 {
				pct := float64(written) / float64(total) * 100
				fmt.Fprintf(output, "\rDownloading ICIS installer... %.0f%%", pct)
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return "", readErr
		}
	}
	println()
	f.Close()

	return destPath, nil
}
