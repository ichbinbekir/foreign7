package data

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

var ActiveLists = make(map[string]bool)
var SelectedList string // Navigasyon sırasında seçilen listeyi taşımak için

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

// GetLists lists dizinindeki tüm .txt dosyalarını döner
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

// LoadWordsFromList Belirli bir dosyadan kelimeleri yükler
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

// LoadActiveWords Tüm aktif listelerdeki kelimeleri (tekil olarak) döner
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

// GetWordOriginMap Aktif listelerdeki kelimelerin hangi listelere ait olduğunu döner
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

// CreateList Yeni bir liste dosyası oluşturur
func CreateList(name string) error {
	if !strings.HasSuffix(name, ".txt") {
		name += ".txt"
	}
	path := filepath.Join(GetDataDir(), name)
	
	// Dosya zaten varsa hata dön
	if _, err := os.Stat(path); err == nil {
		return os.ErrExist
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}

// SaveWordToList Belirli bir listeye kelime kaydeder
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
