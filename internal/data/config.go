package data

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

var AppConfig Config
var T map[string]string

type Config struct {
	InterfaceLang string `json:"interface_lang"` // "TR" or "EN"
}

func GetConfig() Config {
	path := filepath.Join(GetDataDir(), "settings.json")
	content, err := os.ReadFile(path)
	if err != nil {
		cfg := Config{InterfaceLang: "EN"}
		_ = LoadTranslations(cfg.InterfaceLang)
		return cfg
	}

	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		cfg = Config{InterfaceLang: "EN"}
	}
	_ = LoadTranslations(cfg.InterfaceLang)
	return cfg
}

func SaveConfig(cfg Config) error {
	path := filepath.Join(GetDataDir(), "settings.json")
	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	AppConfig = cfg
	_ = LoadTranslations(cfg.InterfaceLang)
	return os.WriteFile(path, content, 0644)
}

func LoadTranslations(lang string) error {
	path := filepath.Join("assets", "lang", strings.ToLower(lang)+".json")
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	T = make(map[string]string)
	return json.Unmarshal(content, &T)
}
