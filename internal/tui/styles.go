package tui

import "github.com/charmbracelet/lipgloss"

var (
	clrGreen       = lipgloss.Color("#00ff41")
	clrGreenMid    = lipgloss.Color("#00cc33")
	clrGreenDark   = lipgloss.Color("#005c1a")
	clrGreenSelect = lipgloss.Color("#002d0e")
	clrRed         = lipgloss.Color("#ff4444")
	clrYellow      = lipgloss.Color("#ffd700")
	clrLightGray   = lipgloss.Color("#707070")

	styleTitleBox = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(clrGreenMid).
			Foreground(clrGreen).
			Bold(true).
			Align(lipgloss.Center)

	stylePanelFocused = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(clrGreen)

	stylePanelHeading = lipgloss.NewStyle().
				Foreground(clrGreenMid).
				Bold(true)

	styleColHeader = lipgloss.NewStyle().
			Foreground(clrGreenDark).
			Bold(true)

	styleSelectedRow = lipgloss.NewStyle().
				Background(clrGreenSelect).
				Foreground(clrGreen).
				Bold(true)

	styleIconActive   = lipgloss.NewStyle().Foreground(clrGreen).Bold(true)
	styleIconInactive = lipgloss.NewStyle().Foreground(clrLightGray)
	styleIconPending  = lipgloss.NewStyle().Foreground(clrYellow)

	styleStateActive   = lipgloss.NewStyle().Foreground(clrGreen)
	styleStatePending  = lipgloss.NewStyle().Foreground(clrYellow)
	styleStateInactive = lipgloss.NewStyle().Foreground(clrLightGray)
	styleStateFailed   = lipgloss.NewStyle().Foreground(clrRed).Bold(true)

	styleMuted = lipgloss.NewStyle().Foreground(clrLightGray)

	styleStatusBar     = lipgloss.NewStyle().Background(clrGreenSelect).Foreground(clrGreenMid).Padding(0, 1)
	styleStatusBracket = lipgloss.NewStyle().Foreground(clrGreenDark)
	styleStatusKey     = lipgloss.NewStyle().Foreground(clrGreen).Bold(true)
)

func portIcon(proto string, hasPID bool) string {
	if !hasPID {
		return styleIconInactive.Render("○")
	}
	if proto == "UDP" {
		return styleIconPending.Render("◆")
	}
	return styleIconActive.Render("✓")
}

func protoStyle(proto string) lipgloss.Style {
	if proto == "UDP" {
		return styleStatePending
	}
	return styleStateActive
}
