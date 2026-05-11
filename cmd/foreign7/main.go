package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/forign7/internal/data"
	"github.com/ichbinbekir/forign7/internal/model"
	"github.com/ichbinbekir/tearouter"
	"github.com/ollama/ollama/api"
)

func main() {
	data.AppConfig = data.GetConfig()
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	routes := []tearouter.Route{
		{
			Path: "/",
			Builder: func() tea.Model {
				return model.NewMenuModel()
			},
			Children: []tearouter.Route{
				{
					Path: "test",
					Builder: func() tea.Model {
						words := data.LoadActiveWords()
						return model.NewTestModel(client, words)
					},
				},
				{
					Path: "sentence-test",
					Builder: func() tea.Model {
						words := data.LoadActiveWords()
						return model.NewSentenceTestModel(client, words)
					},
				},
				{
					Path: "settings",
					Builder: func() tea.Model {
						return model.NewSettingsModel()
					},
				},
				{
					Path: "lists",
					Builder: func() tea.Model {
						return model.NewListSelectModel()
					},
					Children: []tearouter.Route{
						{
							Path: "create",
							Builder: func() tea.Model {
								return model.NewCreateListModel()
							},
						},
						{
							Path: "import",
							Builder: func() tea.Model {
								return model.NewImportListModel()
							},
						},
						{
							Path: "manage",
							Builder: func() tea.Model {
								return model.NewManagerModel(client, data.SelectedList)
							},
						},
					},
				},
			},
		},
	}

	routerModel := tearouter.Model{
		InitialRoute: "/",
		Routes:       routes,
	}

	p := tea.NewProgram(routerModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error occurred: %v", err)
		os.Exit(1)
	}
}
