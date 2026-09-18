package extractor

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ProgressFunc func(file string, current, total int)

func Extract(archivePath string, destDir string, onProgress ProgressFunc) ([]string, error) {
	ext := strings.ToLower(filepath.Ext(archivePath))

	switch ext {
	case ".zip":
		return extractZip(archivePath, destDir, onProgress)
	case ".7z":
		return extract7z(archivePath, destDir, onProgress)
	case ".rar":
		return extractRar(archivePath, destDir, onProgress)
	default:
		return nil, fmt.Errorf("unsupported archive format: %s", ext)
	}
}

func extractZip(archivePath string, destDir string, onProgress ProgressFunc) ([]string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	var extractedFiles []string
	total := len(r.File)

	for i, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		if !strings.HasPrefix(filepath.Clean(fpath), filepath.Clean(destDir)+string(os.PathSeparator)) {
			return nil, fmt.Errorf("zip entry escapes target dir: %s", f.Name)
		}

		if onProgress != nil {
			onProgress(f.Name, i+1, total)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return nil, err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return nil, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return nil, err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, fpath)
	}

	return extractedFiles, nil
}

func extract7z(archivePath string, destDir string, onProgress ProgressFunc) ([]string, error) {
	_ = onProgress
	return nil, fmt.Errorf("7z extraction not yet implemented - please install 7z or use zip format")
}

func extractRar(archivePath string, destDir string, onProgress ProgressFunc) ([]string, error) {
	_ = onProgress
	return nil, fmt.Errorf("rar extraction not yet implemented - please install unrar or use zip format")
}
