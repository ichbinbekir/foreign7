# TeaRouter

**A powerful and simple router for Bubble Tea applications, inspired by GoRouter.**

TeaRouter is designed to simplify page (model) management and navigation in complex applications developed with the `bubbletea` TUI framework. It brings the core principles of Flutter's `gorouter` package to the TUI world.

## Features

- **Stack-Based Navigation**: Easily switch between pages with `Push` and `Pop` operations.
- **Declarative Routing**: Define your routes in a clean and readable way.
- **State-Resetting Navigation**: Navigate to a new page by clearing the navigation history using the `Go` method.
- **Page Replacement**: Replace the current page with a new one without removing it from the stack using `Replace`.
- **Middleware Support**: Intercept route transitions to add middleware logic like authentication or logging.

## Installation

```bash
go get github.com/ichbinbekir/tearouter
```

## Quick Start

Below is a basic example of `tearouter` usage, switching between two pages (`home` and `settings`).

```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/your-username/tearouter"
)

// --- Our Page Models ---

// HomePageModel
type HomePageModel struct{}

func (m HomePageModel) Init() tea.Cmd { return nil }
func (m HomePageModel) View() string {
	return "Home Page\n\nPress 's' to go to the settings page.\nPress 'q' to quit."
}
func (m HomePageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			// Push the settings page onto the stack
			return m, tearouter.Redirect(tearouter.Push, "/settings")
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// SettingsPageModel
type SettingsPageModel struct{}

func (m SettingsPageModel) Init() tea.Cmd { return nil }
func (m SettingsPageModel) View() string {
	return "Settings Page\n\nPress 'b' to go back."
}
func (m SettingsPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			// Pop the current page from the stack
			return m, tearouter.Redirect(tearouter.Pop)
		}
	}
	return m, nil
}

// --- Main Application ---

func main() {
	// Define routes

	outes := []tearouter.Route{
		{
			Path:    "/",
			Builder: func() tea.Model { return HomePageModel{} },
		},
		{
			Path:    "/settings",
			Builder: func() tea.Model { return SettingsPageModel{} },
		},
	}

	// Create the router model

	outerModel := tearouter.Model{
		InitialRoute: "/",
		Routes:       routes,
	}

	p := tea.NewProgram(routerModel)
	if err := p.Start(); err != nil {
		fmt.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}
}
```

## Navigation Methods

Navigation is triggered by the `tearouter.Redirect` command.

- `tearouter.Go`: Clears the navigation stack and navigates to the specified target. It's not possible to go back. Typically used for situations like redirecting to the main page after login.
  ```go
  return m, tearouter.Redirect(tearouter.Go, "/home")
  ```

- `tearouter.Push`: Adds a new page onto the current stack. The user can go back.
  ```go
  return m, tearouter.Redirect(tearouter.Push, "/profile")
  ```

- `tearouter.Replace`: Replaces the current (topmost) page on the stack with a new one. The stack size does not change.
  ```go
  return m, tearouter.Redirect(tearouter.Replace, "/profile/edit")
  ```

- `tearouter.Pop`: Removes the topmost page from the stack and returns to the previous one. Returns an error if there is only one page on the stack.
  ```go
  return m, tearouter.Redirect(tearouter.Pop)
  ```

## Middleware Usage

Middleware is a function that processes every navigation request. You can perform logging or prevent a user from accessing a page they are not authorized for.

If the middleware returns `""` (an empty string), navigation proceeds as normal. If it returns a new path, the user is redirected to that path.

```go
// Example: Auth Middleware
var isAuthenticated = false

func authMiddleware(targetPath string) (newPath string) {
	// Don't interfere with navigation to or from the login page
	if targetPath == "/login" || !isAuthenticated {
		return "" // Proceed
	}

	if !isAuthenticated {
		// If the user is not logged in and tries to access a protected page,
		// redirect them to the login page.
		return "/login"
	}

	return "" // Logged in, proceed
}

func main() {
	routerModel := tearouter.Model{
		// ...
		Middleware: authMiddleware,
	}
	// ...
}
```