package data

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var ActiveLists = make(map[string]bool)
var SelectedList string // Used to carry the selected list during navigation
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

type Prompt struct {
	System string `json:"system"`
	User   string `json:"user"`
}

// LoadPrompt loads a prompt from the assets/prompts.json and executes it with data
func LoadPrompt(key string, data any) (Prompt, error) {
	content, err := os.ReadFile("assets/prompts.json")
	if err != nil {
		return Prompt{}, err
	}

	var allPrompts map[string]Prompt
	if err := json.Unmarshal(content, &allPrompts); err != nil {
		return Prompt{}, err
	}

	p, ok := allPrompts[key]
	if !ok {
		return Prompt{}, os.ErrNotExist
	}

	// Execute user template
	tmpl, err := template.New("user").Parse(p.User)
	if err != nil {
		return p, err
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return p, err
	}

	p.User = tpl.String()
	return p, nil
}

func GetDataDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	dir := filepath.Join(cacheDir, "foreign7")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

// ExportList Copies a list from cache to a local destination
func ExportList(filename string, destPath string) error {
	srcPath := filepath.Join(GetDataDir(), filename)
	input, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return os.WriteFile(destPath, input, 0644)
}

// ImportList Copies a local list to the cache directory
func ImportList(srcPath string) error {
	filename := filepath.Base(srcPath)
	if !strings.HasSuffix(filename, ".txt") {
		filename += ".txt"
	}
	destPath := filepath.Join(GetDataDir(), filename)
	
	input, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return os.WriteFile(destPath, input, 0644)
}

// GetLists returns all .txt files in the lists directory
func GetLists() ([]string, error) {
	dir := GetDataDir()
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var lists []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".txt" {
			lists = append(lists, file.Name())
		}
	}
	return lists, nil
}

// LoadWordsFromList loads words from a specific file
func LoadWordsFromList(filename string) ([]string, error) {
	path := filepath.Join(GetDataDir(), filename)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w := strings.TrimSpace(scanner.Text())
		if w != "" {
			words = append(words, w)
		}
	}
	return words, nil
}

// LoadActiveWords returns all words from active lists (uniquely)
func LoadActiveWords() []string {
	wordMap := make(map[string]bool)
	lists, _ := GetLists()
	
	for _, list := range lists {
		if ActiveLists[list] {
			words, _ := LoadWordsFromList(list)
			for _, w := range words {
				wordMap[strings.ToLower(w)] = true
			}
		}
	}

	var result []string
	for w := range wordMap {
		result = append(result, w)
	}
	return result
}

// GetWordOriginMap returns which lists the active words belong to
func GetWordOriginMap() map[string][]string {
	originMap := make(map[string][]string)
	lists, _ := GetLists()

	for _, listName := range lists {
		if ActiveLists[listName] {
			words, _ := LoadWordsFromList(listName)
			for _, w := range words {
				lowerW := strings.ToLower(w)
				originMap[lowerW] = append(originMap[lowerW], listName)
			}
		}
	}
	return originMap
}

// CreateList creates a new list file
func CreateList(name string) error {
	if !strings.HasSuffix(name, ".txt") {
		name += ".txt"
	}
	path := filepath.Join(GetDataDir(), name)
	
	// Return error if file already exists
	if _, err := os.Stat(path); err == nil {
		return os.ErrExist
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}

// DeleteList deletes a specific list file
func DeleteList(name string) error {
	path := filepath.Join(GetDataDir(), name)
	return os.Remove(path)
}

// SaveWordToList saves a word to a specific list
func SaveWordToList(filename, word string) error {
	path := filepath.Join(GetDataDir(), filename)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(word + "\n")
	return err
}
