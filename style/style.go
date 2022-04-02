package style

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	cross    = lipgloss.AdaptiveColor{Light: "#CC241D", Dark: "#FB4934"}
	text     = lipgloss.AdaptiveColor{Light: "#3C3836", Dark: "#EBDBB2"}
	selected = lipgloss.AdaptiveColor{Light: "#928374", Dark: "#928374"}

	checkedIcon = lipgloss.NewStyle().SetString("✗").
			Foreground(cross).
			PaddingRight(1).
			String()
)

func Branch(name string) string {
	return lipgloss.NewStyle().
		Foreground(text).
		Render(name)
}

func BranchSelect(name string) string {
	return checkedIcon + lipgloss.NewStyle().
		Strikethrough(true).
		Foreground(selected).
		Render(name)
}
