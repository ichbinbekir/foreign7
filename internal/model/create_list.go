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

type CreateListModel struct {
	textInput textinput.Model
	err       error
}

func NewCreateListModel() CreateListModel {
	ti := textinput.New()
	ti.Placeholder = "Liste adı (örn: verbs, daily-phrases)..."
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 40

	return CreateListModel{
		textInput: ti,
	}
}

func (m CreateListModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CreateListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tearouter.Redirect(tearouter.Pop)
		case "enter":
			name := strings.TrimSpace(m.textInput.Value())
			if name == "" {
				return m, nil
			}

			if !strings.HasSuffix(name, ".txt") {
				name += ".txt"
			}

			err := data.CreateList(name)
			if err != nil {
				if os.IsExist(err) {
					m.err = fmt.Errorf("bu isimde bir liste zaten var")
				} else {
					m.err = err
				}
				return m, nil
			}

			// Listeyi aktif et ve yönlendir
			data.ActiveLists[name] = true
			data.SelectedList = name
			// Replace kullanıyoruz ki geri gelindiğinde tekrar isim sorma ekranına düşmeyelim
			return m, tearouter.Redirect(tearouter.Replace, "/manage")
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m CreateListModel) View() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("➕ Yeni Liste Oluştur") + "\n\n")
	sb.WriteString("Oluşturulacak listenin adını yazın:\n\n")
	sb.WriteString(m.textInput.View() + "\n")

	if m.err != nil {
		sb.WriteString("\n" + errorStyle.Render("Hata: "+m.err.Error()) + "\n")
	}

	sb.WriteString("\n(Onaylamak için Enter, Vazgeçmek için Esc)")
	return docStyle.Render(sb.String())
}
