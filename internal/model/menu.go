package model

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ichbinbekir/forign7/internal/data"
	"github.com/ichbinbekir/tearouter"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	docStyle          = lipgloss.NewStyle().Margin(1, 2)
)

type item struct {
	key string
}

func (i item) FilterValue() string { return data.T[i.key] }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, data.T[i.key])
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, ""))
		}
	}
	fmt.Fprint(w, fn(str))
}

type MenuModel struct {
	list list.Model
}

func NewMenuModel() MenuModel {
	items := []list.Item{
		item{key: "menu_test_meaning"},
		item{key: "menu_test_sentence"},
		item{key: "menu_library"},
		item{key: "menu_settings"},
		item{key: "menu_exit"},
	}

	l := list.New(items, itemDelegate{}, 0, 0)
	l.Title = data.T["menu_title"]
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	return MenuModel{
		list: l,
	}
}

type checkListsMsg struct {
	empty bool
}

func (m MenuModel) Init() tea.Cmd {
	return func() tea.Msg {
		lists, _ := data.GetLists()
		return checkListsMsg{empty: len(lists) == 0}
	}
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case checkListsMsg:
		if msg.empty {
			return m, tearouter.Redirect(tearouter.Push, "/create-list")
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-4, msg.Height-4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if !ok {
				return m, nil
			}
			switch i.key {
			case "menu_test_meaning":
				return m, tearouter.Redirect(tearouter.Push, "/test")
			case "menu_test_sentence":
				return m, tearouter.Redirect(tearouter.Push, "/sentence-test")
			case "menu_library":
				return m, tearouter.Redirect(tearouter.Push, "/lists")
			case "menu_settings":
				return m, tearouter.Redirect(tearouter.Push, "/settings")
			case "menu_exit":
				return m, tea.Quit
			}
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MenuModel) View() string {
	m.list.Title = data.T["menu_title"]
	return docStyle.Render(m.list.View())
}
