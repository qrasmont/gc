package main

import (
	"fmt"
	"io"
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

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := ""
	if i.selected {

		str = style.CheckedIcon + style.BranchSelect().Render(i.name)
	} else {
		str = style.CurrentBranch().Render(i.name)
	}

	fn := func(s string) string {
		if index == m.Index() {
			return "⮞ " + s
		}
		return "  " + s
	}

	fmt.Fprintf(w, fn(str))
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Delete key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func newKeyMaps() *keyMap {
	return &keyMap{
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
	cursor   int
	selected map[int]struct{}
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

func initialModel() model {
	branches, err := Get()
	if err != nil {
		panic("oops")
	}

	var keys = newKeyMaps()
	l := list.New(branches, itemDelegate{}, 0, 0)
	l.Title = "Branches:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Select, keys.Delete}
	}

	return model{
		keys:     keys,
		branches: branches,
		selected: make(map[int]struct{}),
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

		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.branches)-1 {
				m.cursor++
			}

		case key.Matches(msg, m.keys.Select):
			cc := m.branches[m.cursor].(item).selected
			text := m.branches[m.cursor].(item).name
			m.branches[m.cursor] = item{name: text, selected: !cc}
			m.list.SetItem(m.cursor, m.branches[m.cursor])

		case key.Matches(msg, m.keys.Delete):
			selection := getSelectedList(m)
			Del(selection)

			branches, err := Get()
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
