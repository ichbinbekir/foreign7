package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/forign7/internal/data"
	"github.com/ichbinbekir/tearouter"
)

type listItem struct {
	name      string
	isGlobal  bool
	isCreator bool
	isImporter bool
}

func (i listItem) Title() string {
	if i.isCreator {
		return data.T["lib_create"]
	}
	if i.isImporter {
		return data.T["lib_import"]
	}
	if i.isGlobal {
		return data.T["lib_global"]
	}
	tick := "[ ]"
	if data.ActiveLists[i.name] {
		tick = "[✓]"
	}
	return fmt.Sprintf("%s %s", tick, i.name)
}

func (i listItem) Description() string {
	if i.isCreator {
		return data.T["lib_create_desc"]
	}
	if i.isImporter {
		return data.T["lib_import_desc"]
	}
	if i.isGlobal {
		return data.T["lib_global_desc"]
	}
	words, _ := data.LoadWordsFromList(i.name)
	return fmt.Sprintf(data.T["lib_words_count"], len(words))
}

func (i listItem) FilterValue() string { return i.name }

type ListSelectModel struct {
	list    list.Model
	message string
}

func NewListSelectModel() ListSelectModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = data.T["lib_title"]
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle

	m := ListSelectModel{list: l}
	m.refreshList()
	return m
}

func (m *ListSelectModel) refreshList() {
	files, _ := data.GetLists()
	items := []list.Item{
		listItem{name: "Global", isGlobal: true},
		listItem{name: "Create", isCreator: true},
		listItem{name: "Import", isImporter: true},
	}
	for _, f := range files {
		items = append(items, listItem{name: f})
	}
	m.list.SetItems(items)
}

func (m ListSelectModel) Init() tea.Cmd { return nil }

func (m ListSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-4, msg.Height-6)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tearouter.Redirect(tearouter.Pop)
		case "e":
			i, ok := m.list.SelectedItem().(listItem)
			if ok && !i.isGlobal && !i.isCreator && !i.isImporter {
				dest := i.name + ".exported.txt"
				err := data.ExportList(i.name, dest)
				if err != nil {
					m.message = errorStyle.Render(fmt.Sprintf(data.T["lib_export_error"], err.Error()))
				} else {
					m.message = successStyle.Render(fmt.Sprintf(data.T["lib_export_success"], dest))
				}
			}
			return m, nil
		case "x":
			i, ok := m.list.SelectedItem().(listItem)
			if ok && !i.isGlobal && !i.isCreator && !i.isImporter {
				err := data.DeleteList(i.name)
				if err != nil {
					m.message = errorStyle.Render(fmt.Sprintf(data.T["lib_delete_error"], err.Error()))
				} else {
					delete(data.ActiveLists, i.name)
					m.message = successStyle.Render(fmt.Sprintf(data.T["lib_delete_success"], i.name))
					m.refreshList()
					
					// If no lists are left, redirect to create
					lists, _ := data.GetLists()
					if len(lists) == 0 {
						return m, tearouter.Redirect(tearouter.Replace, "/lists/create")
					}
				}
			}
			return m, nil
		case " ":
			// Toggle active state
			i, ok := m.list.SelectedItem().(listItem)
			if ok && !i.isGlobal && !i.isCreator && !i.isImporter {
				data.ActiveLists[i.name] = !data.ActiveLists[i.name]
				m.refreshList()
			}
			return m, nil
		case "enter":
			i, ok := m.list.SelectedItem().(listItem)
			if ok {
				if i.isCreator {
					return m, tearouter.Redirect(tearouter.Push, "/lists/create")
				}
				if i.isImporter {
					return m, tearouter.Redirect(tearouter.Push, "/lists/import")
				}
				if i.isGlobal {
					data.SelectedList = "" // Global mode
				} else {
					data.SelectedList = i.name // List mode
				}
				return m, tearouter.Redirect(tearouter.Push, "/lists/manage")
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListSelectModel) View() string {
	view := m.list.View()
	if m.message != "" {
		view += "\n" + m.message
	}
	return docStyle.Render(view)
}
