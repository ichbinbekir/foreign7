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
		return "➕ Yeni Liste Oluştur"
	}
	if i.isImporter {
		return "📥 Liste İçe Aktar (.txt)"
	}
	if i.isGlobal {
		return "🌐 Global: Şu an seçili olanlar"
	}
	tick := "[ ]"
	if data.ActiveLists[i.name] {
		tick = "[✓]"
	}
	return fmt.Sprintf("%s %s", tick, i.name)
}

func (i listItem) Description() string {
	if i.isCreator {
		return "Kelime eklemek için yeni bir kategori oluşturun"
	}
	if i.isImporter {
		return "Bilgisayarınızdaki bir .txt dosyasını kütüphaneye kopyalayın"
	}
	if i.isGlobal {
		return "Tüm aktif listelerde arama yapın (Ekleme kapalı)"
	}
	words, _ := data.LoadWordsFromList(i.name)
	return fmt.Sprintf("%d kelime (Export için 'e' tuşuna basın)", len(words))
}

func (i listItem) FilterValue() string { return i.name }

type ListSelectModel struct {
	list    list.Model
	message string
}

func NewListSelectModel() ListSelectModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Kütüphane Yönetimi"
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
					m.message = errorStyle.Render("Export hatası: " + err.Error())
				} else {
					m.message = successStyle.Render("Liste dışa aktarıldı: " + dest)
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
					return m, tearouter.Redirect(tearouter.Push, "/create-list")
				}
				if i.isImporter {
					return m, tearouter.Redirect(tearouter.Push, "/import")
				}
				if i.isGlobal {
					data.SelectedList = "" // Global mod
				} else {
					data.SelectedList = i.name // Liste modu
				}
				return m, tearouter.Redirect(tearouter.Push, "/manage")
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
