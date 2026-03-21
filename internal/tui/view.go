package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"laxy-ports/internal/models"
)

var statusKeys = [][2]string{
	{"↑↓/jk", "navigate"},
	{"K", "kill"},
	{"R", "refresh"},
	{"Q", "quit"},
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	if m.loading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			styleTitleBox.Render("  LAXY-PORTS  (TUI)  "),
		)
	}
	contentH := max(4, m.height-titleBarLines-statusBarLines)
	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderTitleBar(),
		m.renderPortList(m.width, contentH),
		m.renderStatusBar(),
	)
}

func (m Model) renderTitleBar() string {
	return styleTitleBox.Width(max(1, m.width-2)).Render("LAXY-PORTS  (TUI)")
}

func (m Model) renderStatusBar() string {
	parts := make([]string, 0, len(statusKeys))
	for _, kv := range statusKeys {
		parts = append(parts,
			styleStatusBracket.Render("[") +
				styleStatusKey.Render(kv[0]) +
				styleStatusBracket.Render("]") +
				"  " + kv[1],
		)
	}
	bar := strings.Join(parts, "   ")
	if m.statusMsg != "" {
		bar += "   " + styleMuted.Render("│  "+m.statusMsg)
	}
	return styleStatusBar.Width(m.width).Render(bar)
}

func (m Model) renderPortList(width, height int) string {
	iW, iH := width-2, height-2
	if iW < 10 || iH < 4 {
		return stylePanelFocused.Width(max(2, iW)).Height(max(2, iH)).Render("")
	}

	procW := max(8, iW-portColW-protoColW-pidColW-5)

	title := stylePanelHeading.Render("─ OPEN PORTS ") + styleMuted.Render(fmt.Sprintf("(%d)", len(m.ports)))

	hdr := styleColHeader.Render("  "+fmt.Sprintf("%-*s", portColW, "PORT")) +
		" " + styleColHeader.Render(fmt.Sprintf("%-*s", protoColW, "PROTO")) +
		" " + styleColHeader.Render(fmt.Sprintf("%-*s", pidColW, "PID")) +
		" " + styleColHeader.Render(fmt.Sprintf("%-*s", procW, "PROCESS"))

	separator := styleColHeader.Render(strings.Repeat("─", iW))

	listH := max(0, iH-listHeaderLines)
	rows := make([]string, 0, listH)
	for i := m.listOffset; i < len(m.ports) && len(rows) < listH; i++ {
		rows = append(rows, m.renderPortRow(i, m.ports[i], iW, procW))
	}
	if len(m.ports) == 0 {
		rows = append(rows, styleMuted.Render("  no open ports detected"))
	}
	for len(rows) < listH {
		rows = append(rows, "")
	}

	lines := append([]string{title, hdr, separator}, rows...)
	return stylePanelFocused.Width(iW).Height(iH).Render(
		lipgloss.JoinVertical(lipgloss.Left, lines...),
	)
}

func (m Model) renderPortRow(idx int, p models.PortEntry, iW, procW int) string {
	pidStr := strconv.Itoa(p.PID)
	if p.PID == 0 {
		pidStr = "-"
	}

	procName := p.Process
	if runes := []rune(procName); len(runes) > procW {
		procName = string(runes[:procW-1]) + "…"
	}
	procPadded := procName + strings.Repeat(" ", max(0, procW-len([]rune(procName))))

	portStr := fmt.Sprintf("%-*s", portColW, strconv.Itoa(int(p.Port)))
	protoStr := fmt.Sprintf("%-*s", protoColW, p.Protocol)
	pidPadded := fmt.Sprintf("%-*s", pidColW, pidStr)

	if idx == m.cursor {
		return styleSelectedRow.Width(iW).Render("▶ " + portStr + " " + protoStr + " " + pidPadded + " " + procPadded)
	}

	processStyle := styleStateActive
	if p.PID == 0 {
		processStyle = styleMuted
	}

	return lipgloss.NewStyle().Width(iW).Render(
		portIcon(p.Protocol, p.PID > 0) + " " + portStr + " " +
			protoStyle(p.Protocol).Render(protoStr) + " " +
			styleMuted.Render(pidPadded) + " " +
			processStyle.Render(procPadded),
	)
}
