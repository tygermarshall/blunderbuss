package main

import (
	"fmt"
	"github.com/tygermarshall/blunderbuss/shared/outline"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the application state
type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	width    int
	heigh    int
}

// Init returns the initial command for the application to run
func (m model) Init() tea.Cmd {
	return nil // Just return nil, no commands to run
}

// Update handles events and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	// Return the updated model and a nil command
	return m, nil
}

// View renders the UI based on the model's state
func (m model) View() string {
	var output strings.Builder
	output.WriteString(outline.Top())
	output.WriteString(outline.Middle())
	output.WriteString(outline.Middle())
	output.WriteString(outline.Middle())
	output.WriteString(outline.Middle())
	output.WriteString(outline.Middle())
	output.WriteString(outline.Middle())
	output.WriteString(outline.Bottom())
	return output.String()
}

func main() {
	p := tea.NewProgram(model{})       // Initialize the program with an empty model
	if _, err := p.Run(); err != nil { // Run the program
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
