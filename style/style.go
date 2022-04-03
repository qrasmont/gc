package style

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	cross    = lipgloss.AdaptiveColor{Light: "#CC241D", Dark: "#FB4934"}
	text     = lipgloss.AdaptiveColor{Light: "#3C3836", Dark: "#EBDBB2"}
	selected = lipgloss.AdaptiveColor{Light: "#928374", Dark: "#928374"}

	CheckedIcon = lipgloss.NewStyle().SetString("✗").
			Foreground(cross).
			PaddingRight(1).
			String()
	App = lipgloss.NewStyle().Padding(1, 2)
)

func CurrentBranch() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(text)
}

func BranchSelect() lipgloss.Style {
	return lipgloss.NewStyle().
		Strikethrough(true).
		Foreground(selected)
}
