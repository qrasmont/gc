package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qrasmont/gc/git"
	"github.com/qrasmont/gc/style"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Delete key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.Delete},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys(" ", "s"),
		key.WithHelp("space/s", "select branch"),
	),
	Delete: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "delete selection"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type model struct {
	keys     keyMap
	help     help.Model
	branches []string
	cursor   int
	selected map[int]struct{}
}

func getSelectedList(m model) []string {
	var selection = []string{}

	for i := range m.selected {
		selection = append(selection, m.branches[i])
	}

	return selection
}

func initialModel() model {
	branches, err := git.Get()
	if err != nil {
		panic("oops")
	}
	return model{
		keys:     keys,
		help:     help.New(),
		branches: branches,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.branches)-1 {
				m.cursor++
			}

		case key.Matches(msg, m.keys.Select):
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		case key.Matches(msg, m.keys.Delete):
			selection := getSelectedList(m)
			git.Del(selection)
			m = initialModel()
		}
	}

	return m, nil
}

func (m model) View() string {
	// Header
	s := "Branches \n\n"

	for i, branch := range m.branches {

		// Set the cursor
		cursor := " "
		if m.cursor == i {
			cursor = "⮞"
		}

		// Set selection
		text := ""
		if _, ok := m.selected[i]; ok {
			text = style.BranchSelect(branch)
		} else {
			text = style.Branch(branch)
		}

		// Render row
		s += fmt.Sprintf("%s %s\n", cursor, text)
	}

	// Footer help
	helpView := m.help.View(m.keys)

	return s + "\n" + helpView
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
