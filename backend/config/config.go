package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type AppConfig struct {
	SteamCMDPath  string `json:"steamcmd_path"`
	ServerPath    string `json:"server_path"`
	DefaultPreset string `json:"default_preset"`
}

var (
	cfg *AppConfig
	once sync.Once
	mu   sync.RWMutex
)

func getDefaultPaths() (steamcmd, server string) {
	if runtime.GOOS == "windows" {
		steamcmd = "C:\\steamcmd"
		server = "C:\\ArmaReforgerServer"
	} else {
		home, _ := os.UserHomeDir()
		steamcmd = filepath.Join(home, "steamcmd")
		server = filepath.Join(home, "arma-reforger-server")
	}
	return
}

func getConfigPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "arsm", "config.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "arsm", "config.json")
}

func Load() *AppConfig {
	once.Do(func() {
		steamcmd, server := getDefaultPaths()
		cfg = &AppConfig{
			SteamCMDPath:  steamcmd,
			ServerPath:    server,
			DefaultPreset: "",
		}
		configPath := getConfigPath()
		data, err := os.ReadFile(configPath)
		if err == nil {
			json.Unmarshal(data, cfg)
		}
	})
	return cfg
}

func Get() *AppConfig {
	mu.RLock()
	defer mu.RUnlock()
	return cfg
}

func Save() error {
	mu.Lock()
	defer mu.Unlock()
	configPath := getConfigPath()
	os.MkdirAll(filepath.Dir(configPath), 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func Update(newCfg *AppConfig) error {
	mu.Lock()
	cfg = newCfg
	mu.Unlock()
	return Save()
}
