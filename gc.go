package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qrasmont/gc/git"
)

type model struct {
	branches []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	branches, err := git.Get()
	if err != nil {
		panic("oops")
	}
	return model{
		branches: branches,

		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.branches)-1 {
				m.cursor++
			}

		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			// Here we'll execute the delete then quit
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	// Header
	s := "Branches \n\n"

	for i, choice := range m.branches {

		// Set the cursor
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		// Set selection
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		// Render row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// Footer
	s += "\nQuit: q or ctrl-c\tSelect: space\tDelete: enter\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
