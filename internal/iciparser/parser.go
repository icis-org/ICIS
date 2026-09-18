package iciparser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ICIFile struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Desc        string   `json:"desc"`
	URL         string   `json:"url"`
	Type        string   `json:"type"`     // exe, zip, rar, 7z, msi
	InstallDir  string   `json:"installDir"` // programfiles, appdata, custom
	CustomDir   string   `json:"customDir"`
	Files       []string `json:"files"`
	Shortcut    string   `json:"shortcut"`
	Shortcuts   []ShortcutEntry `json:"shortcuts"`
	Startup     bool     `json:"startup"`
	Icon        string   `json:"icon"`
	Homepage    string   `json:"homepage"`
}

type ShortcutEntry struct {
	Exe string `json:"exe"`
	Name string `json:"name"`
}

func ParseShortcutEntry(raw string) ShortcutEntry {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ShortcutEntry{}
	}
	parts := strings.SplitN(raw, "=", 2)
	exe := strings.TrimSpace(parts[0])
	if strings.EqualFold(filepath.Ext(exe), ".lnk") {
		exe = strings.TrimSuffix(exe, filepath.Ext(exe))
	}
	name := strings.TrimSuffix(filepath.Base(exe), filepath.Ext(exe))
	if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
		name = strings.TrimSpace(parts[1])
	}
	return ShortcutEntry{Exe: exe, Name: name}
}

func Parse(filePath string) (*ICIFile, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open .ici file: %w", err)
	}
	defer f.Close()

	ici := &ICIFile{
		Type:       "zip",
		InstallDir: "appdata",
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch strings.ToLower(key) {
		case "name":
			ici.Name = value
		case "version":
			ici.Version = value
		case "desc":
			ici.Desc = value
		case "url":
			ici.URL = value
		case "type":
			ici.Type = strings.ToLower(value)
		case "install_dir":
			ici.InstallDir = strings.ToLower(value)
		case "custom_dir":
			ici.CustomDir = value
		case "files":
			for _, f := range strings.Split(value, ",") {
				f = strings.TrimSpace(f)
				if f != "" {
					ici.Files = append(ici.Files, f)
				}
			}
		case "shortcut":
			ici.Shortcut = value
		case "shortcuts":
			for _, s := range strings.Split(value, ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					ici.Shortcuts = append(ici.Shortcuts, ParseShortcutEntry(s))
				}
			}
		case "startup":
			ici.Startup = strings.ToLower(value) == "true"
		case "icon":
			ici.Icon = value
		case "homepage":
			ici.Homepage = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .ici file: %w", err)
	}

	if ici.Name == "" {
		return nil, fmt.Errorf("missing required field: name")
	}
	if ici.URL == "" {
		return nil, fmt.Errorf("missing required field: url")
	}

	return ici, nil
}

func (ici *ICIFile) ResolveInstallPath() (string, error) {
	switch ici.InstallDir {
	case "programfiles":
		if runtime.GOOS == "windows" {
			return filepath.Join(os.Getenv("ProgramFiles"), "ICIS", ici.Name), nil
		}
		return filepath.Join("/opt", "icis", ici.Name), nil
	case "appdata":
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(configDir, "ICIS", "apps", ici.Name), nil
	case "custom":
		if ici.CustomDir == "" {
			return "", fmt.Errorf("custom_dir is required when install_dir is custom")
		}
		return filepath.Join(ici.CustomDir, ici.Name), nil
	default:
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(configDir, "ICIS", "apps", ici.Name), nil
	}
}

func ParseString(content string) (*ICIFile, error) {
	ici := &ICIFile{
		Type:       "zip",
		InstallDir: "appdata",
	}

	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch strings.ToLower(key) {
		case "name":
			ici.Name = value
		case "version":
			ici.Version = value
		case "desc":
			ici.Desc = value
		case "url":
			ici.URL = value
		case "type":
			ici.Type = strings.ToLower(value)
		case "install_dir":
			ici.InstallDir = strings.ToLower(value)
		case "custom_dir":
			ici.CustomDir = value
		case "files":
			for _, f := range strings.Split(value, ",") {
				f = strings.TrimSpace(f)
				if f != "" {
					ici.Files = append(ici.Files, f)
				}
			}
		case "shortcut":
			ici.Shortcut = value
		case "shortcuts":
			for _, s := range strings.Split(value, ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					ici.Shortcuts = append(ici.Shortcuts, ParseShortcutEntry(s))
				}
			}
		case "startup":
			ici.Startup = strings.ToLower(value) == "true"
		case "icon":
			ici.Icon = value
		case "homepage":
			ici.Homepage = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .ici content: %w", err)
	}

	if ici.Name == "" {
		return nil, fmt.Errorf("missing required field: name")
	}
	if ici.URL == "" {
		return nil, fmt.Errorf("missing required field: url")
	}

	return ici, nil
}
