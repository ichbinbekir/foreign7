package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ichbinbekir/forign7/internal/data"
	"github.com/ichbinbekir/tearouter"
	"github.com/ollama/ollama/api"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	listTagStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
)

type ollamaValidateMsg struct {
	word   string
	exists bool
	note   string
}

type ManagerModel struct {
	textInput      textinput.Model
	feedback       string
	isValidating   bool
	suggestions    []string
	targetList     string // Boş ise "Global Arama" modu
	allActiveWords []string
	wordOrigins    map[string][]string
	ollama         *api.Client
}

func NewManagerModel(ollama *api.Client, targetList string) ManagerModel {
	ti := textinput.New()
	ti.Placeholder = "Aramak/Eklemek istediğiniz kelime..."
	if targetList == "" {
		ti.Placeholder = "Tüm aktif listelerde ara..."
	}
	ti.Focus()

	return ManagerModel{
		textInput:      ti,
		targetList:     targetList,
		allActiveWords: data.LoadActiveWords(),
		wordOrigins:    data.GetWordOriginMap(),
		ollama:         ollama,
	}
}

func (m ManagerModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tearouter.Redirect(tearouter.Pop)
		case "enter":
			if m.targetList == "" {
				m.feedback = "Global modda kelime ekleme yapılamaz. Sadece arama içindir."
				return m, nil
			}
			if m.isValidating {
				return m, nil
			}
			word := strings.TrimSpace(m.textInput.Value())
			if word == "" {
				return m, nil
			}
			// Mükerrer kontrolü (Tüm aktif listelerde ara)
			for _, w := range m.allActiveWords {
				if strings.EqualFold(w, word) {
					m.feedback = errorStyle.Render("Bu kelime zaten aktif listelerinizde var!")
					return m, nil
				}
			}
			m.isValidating = true
			return m, m.validateWord(word)
		}
	case ollamaValidateMsg:
		m.isValidating = false
		if msg.exists {
			m.saveWord(msg.word)
			m.allActiveWords = append(m.allActiveWords, msg.word)
			// Origin map'i güncelle
			lowerW := strings.ToLower(msg.word)
			m.wordOrigins[lowerW] = append(m.wordOrigins[lowerW], m.targetList)
			
			m.feedback = successStyle.Render(fmt.Sprintf("✓ %s eklendi (%s). %s", msg.word, m.targetList, msg.note))
			m.textInput.Reset()
		} else {
			m.feedback = errorStyle.Render(fmt.Sprintf("✗ %s? %s", msg.word, msg.note))
		}
		return m, nil
	case errMsg:
		m.isValidating = false
		m.feedback = errorStyle.Render("Hata: " + msg.Error())
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	// Suggestions
	input := strings.ToLower(m.textInput.Value())
	if input != "" {
		m.suggestions = []string{}
		for _, w := range m.allActiveWords {
			if strings.Contains(strings.ToLower(w), input) {
				m.suggestions = append(m.suggestions, w)
				if len(m.suggestions) > 10 {
					break
				}
			}
		}
	} else {
		m.suggestions = nil
	}

	return m, cmd
}

func (m ManagerModel) View() string {
	var sb strings.Builder
	title := "Global Arama (Görüntüleme)"
	if m.targetList != "" {
		title = fmt.Sprintf("Liste Yönetimi: %s", m.targetList)
	}

	sb.WriteString(titleStyle.Render(title) + "\n\n")
	
	// Toplam kelime sayısı bilgisi
	countInfo := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(fmt.Sprintf("(Toplam %d benzersiz kelime aktif)", len(m.allActiveWords)))
	sb.WriteString(countInfo + "\n\n")

	if m.targetList != "" {
		sb.WriteString("İngilizce kelime yazın. Transgemma kontrol edip bu listeye kaydedecek.\n\n")
	} else {
		sb.WriteString("Seçili olan tüm listelerde arama yapıyorsunuz.\n\n")
	}

	sb.WriteString(m.textInput.View() + "\n")

	if m.isValidating {
		sb.WriteString("\nKelime doğrulanıyor...")
	}

	if len(m.suggestions) > 0 {
		sb.WriteString("\n\nBulunan / Benzer kelimeler:\n")
		for _, s := range m.suggestions {
			origins := m.wordOrigins[strings.ToLower(s)]
			originText := ""
			if len(origins) > 0 {
				originText = listTagStyle.Render(" [" + strings.Join(origins, ", ") + "]")
			}
			sb.WriteString("- " + s + originText + "\n")
		}
	} else if m.textInput.Value() != "" && !m.isValidating {
		sb.WriteString("\n\n(Bu kelime listelerinizde bulunamadı)")
	}

	if m.feedback != "" {
		sb.WriteString("\n\n" + m.feedback)
	}

	sb.WriteString("\n\n(Geri dönmek için Esc)")
	return docStyle.Render(sb.String())
}

func (m ManagerModel) validateWord(word string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		p, err := data.LoadPrompt("word_validation", map[string]string{"Word": word})
		if err != nil {
			return errMsg(err)
		}

		req := &api.ChatRequest{
			Model: "translategemma:latest",
			Messages: []api.Message{
				{Role: "system", Content: p.System},
				{Role: "user", Content: p.User},
			},
			Stream: new(bool),
		}
		var respStr string
		err = m.ollama.Chat(ctx, req, func(r api.ChatResponse) error {
			respStr = r.Message.Content
			return nil
		})
		if err != nil {
			return errMsg(err)
		}

		exists := strings.HasPrefix(strings.ToUpper(respStr), "YES")
		return ollamaValidateMsg{word: word, exists: exists, note: respStr}
	}
}

func (m ManagerModel) saveWord(word string) {
	_ = data.SaveWordToList(m.targetList, word)
}
