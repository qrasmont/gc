package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qrasmont/gc/style"
)

var keepSelect bool = false

type itemDelegate struct{}

func SetSelect(m *list.Model) {
	index := m.Index()
	currentItem := m.SelectedItem()
	m.SetItem(index, item{name: currentItem.(item).name,
		selected: true})
}

func (d itemDelegate) Height() int  { return 1 }
func (d itemDelegate) Spacing() int { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.KeepSelect):
			keepSelect = !keepSelect

			if keepSelect {
                SetSelect(m)
			}

		case
			key.Matches(msg, m.KeyMap.CursorUp),
			key.Matches(msg, m.KeyMap.CursorDown):
			if keepSelect {
                SetSelect(m)
			}
		}
	}

	return nil
}

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
