package model

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/forign7/internal/data"
	"github.com/ichbinbekir/tearouter"
)

type ImportListModel struct {
	textInput textinput.Model
	err       error
	success   string
}

func NewImportListModel() ImportListModel {
	ti := textinput.New()
	ti.Placeholder = data.T["import_placeholder"]
	ti.Focus()
	ti.Width = 60

	return ImportListModel{
		textInput: ti,
	}
}

func (m ImportListModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ImportListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tearouter.Redirect(tearouter.Pop)
		case "enter":
			path := strings.TrimSpace(m.textInput.Value())
			if path == "" {
				return m, nil
			}

			err := data.ImportList(path)
			if err != nil {
				if os.IsNotExist(err) {
					m.err = fmt.Errorf(data.T["import_error_not_found"])
				} else {
					m.err = err
				}
				return m, nil
			}

			m.success = data.T["import_success"]
			m.textInput.Reset()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m ImportListModel) View() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render(data.T["import_title"]) + "\n\n")
	sb.WriteString(data.T["import_prompt"] + "\n\n")
	sb.WriteString(m.textInput.View() + "\n")

	if m.err != nil {
		sb.WriteString("\n" + errorStyle.Render("Hata: "+m.err.Error()) + "\n")
	}
	if m.success != "" {
		sb.WriteString("\n" + successStyle.Render(m.success) + "\n")
	}

	sb.WriteString("\n" + data.T["import_confirm"])
	return docStyle.Render(sb.String())
}
