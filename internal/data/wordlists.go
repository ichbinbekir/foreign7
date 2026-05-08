package data

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

var ActiveLists = make(map[string]bool)
var SelectedList string // Used to carry the selected list during navigation

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

// RemoveWordFromList removes a word from a specific list file
func RemoveWordFromList(filename, word string) error {
	words, err := LoadWordsFromList(filename)
	if err != nil {
		return err
	}

	var sb strings.Builder
	for _, w := range words {
		if !strings.EqualFold(w, word) {
			sb.WriteString(w + "\n")
		}
	}

	path := filepath.Join(GetDataDir(), filename)
	return os.WriteFile(path, []byte(sb.String()), 0644)
}
