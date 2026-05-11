# TeaRouter

**A powerful and simple router for Bubble Tea applications, inspired by GoRouter.**

TeaRouter is designed to simplify page (model) management and navigation in complex applications developed with the `bubbletea` TUI framework. It brings the core principles of Flutter's `gorouter` package to the TUI world.

## Features

- **Hierarchical Sub-Routing**: Define nested routes and automatically build the navigation stack.
- **Stack-Based Navigation**: Easily switch between pages with `Push` and `Pop` operations.
- **Declarative Routing**: Define your routes in a clean and readable way.
- **State-Resetting Navigation**: Navigate to a new page by clearing the navigation history using the `Go` method.
- **Page Replacement**: Replace the current page with a new one without removing it from the stack using `Replace`.
- **Middleware Support**: Intercept route transitions to add middleware logic like authentication or logging.

## Installation

```bash
go get github.com/ichbinbekir/tearouter
```

## Hierarchical Sub-Routing

Inspired by **GoRouter**, TeaRouter supports hierarchical route definitions. When you navigate to a deep path like `/main/settings/profile`, TeaRouter automatically builds the entire stack of parent models. This allows natural `Pop` behavior (going from Profile to Settings, then to Main).

```go
routes := []tearouter.Route{
    {
        Path: "/main",
        Builder: func() tea.Model { return MainModel{} },
        Children: []tearouter.Route{
            {
                Path: "settings", // Relative path: /main/settings
                Builder: func() tea.Model { return SettingsModel{} },
                Children: []tearouter.Route{
                    {
                        Path: "profile", // Relative path: /main/settings/profile
                        Builder: func() tea.Model { return ProfileModel{} },
                    },
                },
            },
        },
    },
}
```

When you call `tearouter.Redirect(tearouter.Go, "/main/settings/profile")`:
1. The stack is cleared.
2. `MainModel`, `SettingsModel`, and `ProfileModel` are created and stacked in order.
3. The user sees the `ProfileModel`.
4. Calling `Pop` will naturally return the user to `SettingsModel`.

## Quick Start

Below is a basic example of `tearouter` usage, switching between pages.

```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/tearouter"
)

// --- Our Page Models ---

type HomePageModel struct{}
func (m HomePageModel) Init() tea.Cmd { return nil }
func (m HomePageModel) View() string { return "Home Page\n\nPress 's' for Settings\nPress 'q' to quit." }
func (m HomePageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "s": return m, tearouter.Redirect(tearouter.Push, "/settings")
		case "q": return m, tea.Quit
		}
	}
	return m, nil
}

type SettingsPageModel struct{}
func (m SettingsPageModel) Init() tea.Cmd { return nil }
func (m SettingsPageModel) View() string { return "Settings Page\n\nPress 'b' to go back." }
func (m SettingsPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "b" {
		return m, tearouter.Redirect(tearouter.Pop)
	}
	return m, nil
}

func main() {
	routes := []tearouter.Route{
		{ Path: "/", Builder: func() tea.Model { return HomePageModel{} } },
		{ Path: "/settings", Builder: func() tea.Model { return SettingsPageModel{} } },
	}

	routerModel := tearouter.Model{
		InitialRoute: "/",
		Routes:       routes,
	}

	if _, err := tea.NewProgram(routerModel).Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
```

## Navigation Methods

Navigation is triggered by the `tearouter.Redirect` command.

- `tearouter.Go`: Builds the full hierarchy of the target path and sets it as the new stack.
- `tearouter.Push`: Adds the full hierarchy of the target path on top of the current stack.
- `tearouter.Replace`: Replaces the current stack with the full hierarchy of the target path.
- `tearouter.Pop`: Removes the topmost page from the stack and returns to the previous one.

## Middleware Usage

Middleware allows you to intercept navigation requests for tasks like authentication.

```go
func authMiddleware(targetPath string) (newPath string) {
	if !isLoggedIn && targetPath != "/login" {
		return "/login"
	}
	return ""
}
```
