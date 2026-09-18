package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ProgressFunc func(downloaded, total int64, percent float64)

func Download(url string, destDir string, onProgress ProgressFunc) (string, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create download directory: %w", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "ICIS/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	fileName := getFileName(resp, url)
	destPath := filepath.Join(destDir, fileName)

	out, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	totalSize := resp.ContentLength

	buf := make([]byte, 32*1024)
	var downloaded int64

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				return "", writeErr
			}
			downloaded += int64(n)

			if onProgress != nil {
				percent := 0.0
				if totalSize > 0 {
					percent = float64(downloaded) / float64(totalSize) * 100
				}
				onProgress(downloaded, totalSize, percent)
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}

	return destPath, nil
}

func getFileName(resp *http.Response, url string) string {
	cd := resp.Header.Get("Content-Disposition")
	if cd != "" {
		for _, part := range splitCD(cd) {
			if len(part) > 10 && part[:10] == "filename=\"" {
				name := part[10:]
				if idx := len(name) - 1; idx >= 0 && name[idx] == '"' {
					name = name[:idx]
				}
				if name != "" {
					return name
				}
			}
		}
	}

	return filepath.Base(url)
}

func splitCD(cd string) []string {
	var parts []string
	current := ""
	for _, c := range cd {
		if c == ';' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	parts = append(parts, current)
	return parts
}
