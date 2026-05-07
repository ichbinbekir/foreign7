package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/forign7/internal/data"
	"github.com/ichbinbekir/tearouter"
)

type settingsItem struct {
	title       string
	description string
	action      func() tea.Cmd
}

func (i settingsItem) Title() string       { return i.title }
func (i settingsItem) Description() string { return i.description }
func (i settingsItem) FilterValue() string { return i.title }

type SettingsModel struct {
	list list.Model
}

func NewSettingsModel() SettingsModel {
	items := []list.Item{
		settingsItem{
			title:       data.T["settings_lang"],
			description: fmt.Sprintf(data.T["settings_lang_desc"], data.AppConfig.InterfaceLang),
			action: func() tea.Cmd {
				cfg := data.AppConfig
				if cfg.InterfaceLang == "TR" {
					cfg.InterfaceLang = "EN"
				} else {
					cfg.InterfaceLang = "TR"
				}
				data.SaveConfig(cfg)
				return nil
			},
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = data.T["settings_title"]
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	return SettingsModel{list: l}
}

func (m SettingsModel) Init() tea.Cmd { return nil }

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-4, msg.Height-4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tearouter.Redirect(tearouter.Pop)
		case "enter":
			if i, ok := m.list.SelectedItem().(settingsItem); ok {
				cmd := i.action()
				// Refresh list to show new value with new translations
				return NewSettingsModel(), cmd
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SettingsModel) View() string {
	return docStyle.Render(m.list.View())
}
