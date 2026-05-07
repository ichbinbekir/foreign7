package data

import (
	"bytes"
	"encoding/json"
	"os"
	"text/template"
)

type Prompt struct {
	System string `json:"system"`
	User   string `json:"user"`
}

// LoadPrompt loads a prompt from the assets/prompts.json and executes it with data
func LoadPrompt(key string, data any) (Prompt, error) {
	content, err := os.ReadFile("assets/prompts.json")
	if err != nil {
		return Prompt{}, err
	}

	var allPrompts map[string]Prompt
	if err := json.Unmarshal(content, &allPrompts); err != nil {
		return Prompt{}, err
	}

	p, ok := allPrompts[key]
	if !ok {
		return Prompt{}, os.ErrNotExist
	}

	// Execute user template
	tmpl, err := template.New("user").Parse(p.User)
	if err != nil {
		return p, err
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return p, err
	}

	p.User = tpl.String()
	return p, nil
}
