package model

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ichbinbekir/forign7/internal/data"
	"github.com/ichbinbekir/tearouter"
	"github.com/ollama/ollama/api"
)

var (
	wordStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
)

type TestModel struct {
	words          []string
	currentWordIdx int
	textInput      textinput.Model
	feedback       string
	isChecking     bool
	ollama         *api.Client
}

type ollamaCheckMsg string
type nextWordMsg struct{}
type errMsg error

func NewTestModel(ollama *api.Client, words []string) TestModel {
	ti := textinput.New()
	ti.Placeholder = "Anlam..."
	ti.Focus()

	return TestModel{
		words:     words,
		textInput: ti,
		ollama:    ollama,
	}
}

func (m TestModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.isChecking {
		switch msg := msg.(type) {
		case ollamaCheckMsg:
			m.isChecking = false
			m.feedback = string(msg)
			return m, nil
		case errMsg:
			m.isChecking = false
			m.feedback = errorStyle.Render("Hata: " + msg.Error())
			return m, nil
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tearouter.Redirect(tearouter.Pop)
		case "enter":
			if m.feedback != "" {
				return m, func() tea.Msg { return nextWordMsg{} }
			}
			answer := m.textInput.Value()
			if strings.TrimSpace(answer) == "" {
				return m, nil
			}
			m.isChecking = true
			return m, m.checkAnswer(m.words[m.currentWordIdx], answer)
		}
	case nextWordMsg:
		m.feedback = ""
		m.currentWordIdx = (m.currentWordIdx + 1) % len(m.words)
		m.textInput.Reset()
		return m, nil
	case ollamaCheckMsg:
		m.isChecking = false
		m.feedback = string(msg)
		return m, nil
	case errMsg:
		m.isChecking = false
		m.feedback = errorStyle.Render("Hata: " + msg.Error())
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m TestModel) View() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Anlam Tahmini Testi") + "\n\n")
	
	if len(m.words) == 0 {
		return docStyle.Render("Kelimeler yüklenemedi. Liste boş olabilir.")
	}

	sb.WriteString("Soru: " + wordStyle.Render(m.words[m.currentWordIdx]) + "\n\n")
	
	if m.feedback != "" {
		sb.WriteString(m.feedback + "\n\n")
		sb.WriteString("(Sonraki için Enter, Menü için Esc)")
	} else {
		sb.WriteString(m.textInput.View())
		if m.isChecking {
			sb.WriteString("\n\nKontrol ediliyor...")
		}
	}
	return docStyle.Render(sb.String())
}

func (m TestModel) checkAnswer(word, answer string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		p, err := data.LoadPrompt("translation_check", map[string]string{
			"Word":  word,
			"Guess": answer,
		})
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
		return ollamaCheckMsg(respStr)
	}
}
