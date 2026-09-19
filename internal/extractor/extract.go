package extractor

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/nwaples/rardecode"
)

type ProgressFunc func(file string, current, total int)

func Extract(archivePath string, destDir string, onProgress ProgressFunc) ([]string, error) {
	return ExtractAs(archivePath, "", destDir, onProgress)
}

func ExtractAs(archivePath, archiveType, destDir string, onProgress ProgressFunc) ([]string, error) {
	ext := strings.ToLower(strings.TrimSpace(archiveType))
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	if ext == "" {
		ext = strings.ToLower(filepath.Ext(archivePath))
	}

	switch ext {
	case ".zip":
		return extractZip(archivePath, destDir, onProgress)
	case ".7z":
		return extract7z(archivePath, destDir, onProgress)
	case ".rar":
		return extractRar(archivePath, destDir, onProgress)
	case ".exe":
		return extractFile(archivePath, destDir, onProgress)
	default:
		return nil, fmt.Errorf("unsupported archive format: %s", ext)
	}
}

func extractFile(sourcePath, destDir string, onProgress ProgressFunc) ([]string, error) {
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return nil, err
	}
	name := filepath.Base(sourcePath)
	targetPath := filepath.Join(destDir, name)
	if onProgress != nil {
		onProgress(name, 1, 1)
	}
	in, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer in.Close()
	out, err := os.Create(targetPath)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return nil, err
	}
	if err := out.Close(); err != nil {
		return nil, err
	}
	return []string{targetPath}, nil
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
	r, err := sevenzip.OpenReader(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open 7z: %w", err)
	}
	defer r.Close()

	var extractedFiles []string
	total := len(r.File)
	for i, f := range r.File {
		fpath, err := safeTargetPath(destDir, f.Name)
		if err != nil {
			return nil, fmt.Errorf("7z entry %q: %w", f.Name, err)
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
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		out, err := os.Create(fpath)
		if err != nil {
			rc.Close()
			return nil, err
		}
		_, copyErr := io.Copy(out, rc)
		closeErr := out.Close()
		rc.Close()
		if copyErr != nil {
			return nil, copyErr
		}
		if closeErr != nil {
			return nil, closeErr
		}
		extractedFiles = append(extractedFiles, fpath)
	}
	return extractedFiles, nil
}

func extractRar(archivePath string, destDir string, onProgress ProgressFunc) ([]string, error) {
	r, err := rardecode.OpenReader(archivePath, "")
	if err != nil {
		return nil, fmt.Errorf("failed to open rar: %w", err)
	}
	defer r.Close()

	var extractedFiles []string
	var current int
	for {
		header, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read rar: %w", err)
		}
		current++
		fpath, err := safeTargetPath(destDir, header.Name)
		if err != nil {
			return nil, fmt.Errorf("rar entry %q: %w", header.Name, err)
		}
		if onProgress != nil {
			onProgress(header.Name, current, 0)
		}
		if header.IsDir {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return nil, err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return nil, err
		}
		out, err := os.Create(fpath)
		if err != nil {
			return nil, err
		}
		_, copyErr := io.Copy(out, r)
		closeErr := out.Close()
		if copyErr != nil {
			return nil, copyErr
		}
		if closeErr != nil {
			return nil, closeErr
		}
		extractedFiles = append(extractedFiles, fpath)
	}
	return extractedFiles, nil
}

func safeTargetPath(destDir, name string) (string, error) {
	target := filepath.Join(destDir, filepath.FromSlash(name))
	cleanDest := filepath.Clean(destDir)
	if target != cleanDest && !strings.HasPrefix(target, cleanDest+string(os.PathSeparator)) {
		return "", fmt.Errorf("entry escapes target directory")
	}
	return target, nil
}
