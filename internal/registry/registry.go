package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultRegistryURL = "https://raw.githubusercontent.com/icis-org/registry/main/registry.json"
	MaxIndexSize       = 5 * 1024 * 1024
)

type RegistryIndex struct {
	Name     string            `json:"name"`
	Packages []RegistryPackage `json:"packages"`
}

type RegistryPackage struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Desc     string `json:"desc"`
	ICIURL   string `json:"ici"`
	Homepage string `json:"homepage"`
}

type RegistryResult struct {
	Index *RegistryIndex `json:"index"`
	Error string         `json:"error,omitempty"`
}

func configPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "ICIS", "registry-config.json"), nil
}

type RegistryConfig struct {
	RegistryURL string `json:"registryURL"`
}

func LoadConfig() RegistryConfig {
	cfg := RegistryConfig{RegistryURL: DefaultRegistryURL}
	path, err := configPath()
	if err != nil {
		return cfg
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	json.Unmarshal(data, &cfg)
	if cfg.RegistryURL == "" {
		cfg.RegistryURL = DefaultRegistryURL
	}
	return cfg
}

func SaveConfig(cfg RegistryConfig) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func FetchIndex(url string) RegistryResult {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RegistryResult{Error: fmt.Sprintf("request failed: %v", err)}
	}
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req.Header.Set("Pragma", "no-cache")
	resp, err := client.Do(req)
	if err != nil {
		return RegistryResult{Error: fmt.Sprintf("fetch failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RegistryResult{Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxIndexSize))
	if err != nil {
		return RegistryResult{Error: fmt.Sprintf("read failed: %v", err)}
	}

	var idx RegistryIndex
	if err := json.Unmarshal(body, &idx); err != nil {
		return RegistryResult{Error: fmt.Sprintf("parse failed: %v", err)}
	}

	return RegistryResult{Index: &idx}
}

func FetchICI(url string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req.Header.Set("Pragma", "no-cache")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch .ici: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d fetching .ici", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxIndexSize))
	if err != nil {
		return "", fmt.Errorf("failed to read .ici: %w", err)
	}

	return string(body), nil
}
