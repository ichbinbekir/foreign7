package model

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ichbinbekir/forign7/internal/data"
	"github.com/ichbinbekir/tearouter"
	"github.com/ollama/ollama/api"
)

type SentenceTestModel struct {
	words          []string
	currentWordIdx int
	textarea       textarea.Model
	feedback       string
	isChecking     bool
	ollama         *api.Client
}

func NewSentenceTestModel(ollama *api.Client, words []string) SentenceTestModel {
	ta := textarea.New()
	ta.Placeholder = data.T["sentence_placeholder"]
	ta.Focus()
	ta.SetHeight(3)

	return SentenceTestModel{
		words:    words,
		textarea: ta,
		ollama:   ollama,
	}
}

func (m SentenceTestModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m SentenceTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			sentence := m.textarea.Value()
			if strings.TrimSpace(sentence) == "" {
				return m, nil
			}
			m.isChecking = true
			return m, m.checkSentence(m.words[m.currentWordIdx], sentence)
		}
	case nextWordMsg:
		m.feedback = ""
		m.currentWordIdx = (m.currentWordIdx + 1) % len(m.words)
		m.textarea.Reset()
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
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m SentenceTestModel) View() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render(data.T["sentence_title"]) + "\n\n")

	if len(m.words) == 0 {
		return docStyle.Render(data.T["test_error_no_words"])
	}

	sb.WriteString(data.T["sentence_prefix"] + wordStyle.Render(m.words[m.currentWordIdx]) + "\n\n")

	if m.feedback != "" {
		fbStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("2")).
			Padding(1).
			Width(60)

		sb.WriteString(data.T["sentence_your_sentence"] + "\n" + m.textarea.Value() + "\n\n")
		sb.WriteString(fbStyle.Render(m.feedback) + "\n")
		sb.WriteString("\n" + data.T["test_next_info"])
	} else {
		sb.WriteString(m.textarea.View())
		if m.isChecking {
			sb.WriteString("\n\n" + data.T["sentence_analyzing"])
		} else {
			sb.WriteString("\n\n" + data.T["sentence_hint"])
		}
	}
	return docStyle.Render(sb.String())
}

func (m SentenceTestModel) checkSentence(word, sentence string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		p, err := data.LoadPrompt("sentence_check", map[string]string{
			"Word":     word,
			"Sentence": sentence,
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
