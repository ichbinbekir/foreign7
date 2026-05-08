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
	targetList     string // If empty, "Global Search" mode
	allActiveWords []string
	wordOrigins    map[string][]string
	ollama         *api.Client
	selectedIndex  int
}

func NewManagerModel(ollama *api.Client, targetList string) ManagerModel {
	ti := textinput.New()
	ti.Placeholder = data.T["manage_placeholder"]
	if targetList == "" {
		ti.Placeholder = data.T["manage_placeholder_global"]
	}
	ti.Focus()

	return ManagerModel{
		textInput:      ti,
		targetList:     targetList,
		allActiveWords: data.LoadActiveWords(),
		wordOrigins:    data.GetWordOriginMap(),
		ollama:         ollama,
		selectedIndex:  -1,
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
		case "up":
			if len(m.suggestions) > 0 {
				m.selectedIndex--
				if m.selectedIndex < -1 {
					m.selectedIndex = len(m.suggestions) - 1
				}
			}
			return m, nil
		case "down":
			if len(m.suggestions) > 0 {
				m.selectedIndex++
				if m.selectedIndex >= len(m.suggestions) {
					m.selectedIndex = -1
				}
			}
			return m, nil
		case "x":
			if m.selectedIndex >= 0 && m.selectedIndex < len(m.suggestions) {
				wordToDelete := m.suggestions[m.selectedIndex]
				origins := m.wordOrigins[strings.ToLower(wordToDelete)]
				for _, listName := range origins {
					_ = data.RemoveWordFromList(listName, wordToDelete)
				}
				m.feedback = successStyle.Render(fmt.Sprintf(data.T["manage_deleted_success"], wordToDelete))
				// Refresh data
				m.allActiveWords = data.LoadActiveWords()
				m.wordOrigins = data.GetWordOriginMap()
				m.updateSuggestions()
				return m, nil
			}
		case "enter":
			if m.targetList == "" {
				m.feedback = data.T["manage_global_error"]
				return m, nil
			}
			if m.isValidating {
				return m, nil
			}
			word := strings.TrimSpace(m.textInput.Value())
			if word == "" {
				return m, nil
			}
			// Duplication check (Search across all active lists)
			for _, w := range m.allActiveWords {
				if strings.EqualFold(w, word) {
					m.feedback = errorStyle.Render(data.T["manage_exists_error"])
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
			// Update origin map
			lowerW := strings.ToLower(msg.word)
			m.wordOrigins[lowerW] = append(m.wordOrigins[lowerW], m.targetList)
			
			m.feedback = successStyle.Render(fmt.Sprintf(data.T["manage_added_success"], msg.word, m.targetList, msg.note))
			m.textInput.Reset()
			m.updateSuggestions()
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
	m.updateSuggestions()

	return m, cmd
}

func (m *ManagerModel) updateSuggestions() {
	input := strings.ToLower(m.textInput.Value())
	if input != "" {
		newSuggestions := []string{}
		for _, w := range m.allActiveWords {
			if strings.Contains(strings.ToLower(w), input) {
				newSuggestions = append(newSuggestions, w)
				if len(newSuggestions) > 10 {
					break
				}
			}
		}
		m.suggestions = newSuggestions
		if m.selectedIndex >= len(m.suggestions) {
			m.selectedIndex = len(m.suggestions) - 1
		}
	} else {
		m.suggestions = nil
		m.selectedIndex = -1
	}
}

func (m ManagerModel) View() string {
	var sb strings.Builder
	title := data.T["manage_title_global"]
	if m.targetList != "" {
		title = fmt.Sprintf(data.T["manage_title_list"], m.targetList)
	}

	sb.WriteString(titleStyle.Render(title) + "\n\n")
	
	// Total word count info
	countInfo := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(fmt.Sprintf(data.T["manage_total_words"], len(m.allActiveWords)))
	sb.WriteString(countInfo + "\n\n")

	if m.targetList != "" {
		sb.WriteString(data.T["manage_desc_list"])
	} else {
		sb.WriteString(data.T["manage_desc_global"])
	}

	sb.WriteString(m.textInput.View() + "\n")

	if m.isValidating {
		sb.WriteString("\n" + data.T["manage_validating"])
	}

	if len(m.suggestions) > 0 {
		sb.WriteString("\n\n" + data.T["manage_suggestions"] + "\n")
		for i, s := range m.suggestions {
			origins := m.wordOrigins[strings.ToLower(s)]
			originText := ""
			if len(origins) > 0 {
				originText = listTagStyle.Render(" [" + strings.Join(origins, ", ") + "]")
			}
			prefix := "- "
			style := lipgloss.NewStyle()
			if i == m.selectedIndex {
				prefix = "> "
				style = style.Foreground(lipgloss.Color("170")).Bold(true)
			}
			sb.WriteString(style.Render(prefix+s) + originText + "\n")
		}
		if m.selectedIndex >= 0 {
			sb.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(data.T["manage_delete_info"]))
		}
	} else if m.textInput.Value() != "" && !m.isValidating {
		sb.WriteString("\n\n" + data.T["manage_not_found"])
	}

	if m.feedback != "" {
		sb.WriteString("\n\n" + m.feedback)
	}

	sb.WriteString("\n\n" + data.T["manage_back"])
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
