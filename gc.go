package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qrasmont/gc/style"
)

type item struct {
	name     string
	selected bool
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.name }

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	KeepSelect key.Binding
	Select     key.Binding
	Delete     key.Binding
	Help       key.Binding
	Quit       key.Binding
}

func newKeyMaps() *keyMap {
	return &keyMap{
		KeepSelect: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "keep select"),
		),
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
}

type model struct {
	keys     *keyMap
	branches []list.Item
	list     list.Model
}

func getSelectedList(m model) []string {
	var selection = []string{}

	for _, branch := range m.branches {
		if branch.(item).selected {
			selection = append(selection, branch.(item).name)
		}
	}

	return selection
}

var keys = newKeyMaps()

func initialModel() model {
	branches, err := GitGetBranches()
	if err != nil {
		panic("oops")
	}

	l := list.New(branches, itemDelegate{}, 0, 0)
	l.Title = "Branches:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Select, keys.Delete, keys.KeepSelect}
	}

	return model{
		keys:     keys,
		branches: branches,
		list:     l,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Select):
			index := m.list.Index()
			cc := m.branches[index].(item).selected
			text := m.branches[index].(item).name
			m.branches[index] = item{name: text, selected: !cc}
			m.list.SetItem(index, m.branches[index])

		case key.Matches(msg, m.keys.Delete):
			selection := getSelectedList(m)
			GitDelete(selection)

			branches, err := GitGetBranches()
			if err != nil {
				panic("oops")
			}

			m.branches = branches
			m.list.SetItems(m.branches)
		}
	case tea.WindowSizeMsg:
		h, v := style.App.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return style.App.Render(m.list.View())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
