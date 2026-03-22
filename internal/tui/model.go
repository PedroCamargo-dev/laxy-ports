package tui

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"laxy-ports/internal/models"
	"laxy-ports/internal/network"
)

const (
	titleBarLines    = 3
	statusBarLines   = 1
	panelBorderLines = 2
	listHeaderLines  = 3
	portColW         = 6
	protoColW        = 6
	pidColW          = 8
)

type (
	tickMsg     time.Time
	portsMsg    []models.PortEntry
	killDoneMsg struct {
		pid int
		err error
	}
)

type Model struct {
	ports      []models.PortEntry
	cursor     int
	listOffset int
	width      int
	height     int
	statusMsg  string
	loading    bool
}

func New() Model { return Model{loading: true} }

func (m Model) Init() tea.Cmd {
	return tea.Batch(loadPortsCmd(), tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.MouseMsg:
		return m.handleMouse(msg)
	case tickMsg:
		return m, tea.Batch(loadPortsCmd(), tickCmd())
	case portsMsg:
		prev := selectedKey(m.ports, m.cursor)
		m.ports, m.loading = []models.PortEntry(msg), false
		m.cursor = 0
		for i := range m.ports {
			if portEntryKey(m.ports[i]) == prev {
				m.cursor = i
				break
			}
		}
		m.adjustListOffset()
		return m, nil
	case killDoneMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("kill %d: %v", msg.pid, msg.err)
		} else {
			m.statusMsg = fmt.Sprintf("killed PID %d", msg.pid)
		}
		return m, loadPortsCmd()
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.adjustListOffset()
		}
	case "down", "j":
		if m.cursor < len(m.ports)-1 {
			m.cursor++
			m.adjustListOffset()
		}
	case "x":
		if m.cursor < len(m.ports) {
			p := m.ports[m.cursor]
			if p.PID <= 0 {
				m.statusMsg = "no process to kill"
				return m, nil
			}
			return m, killProcess(p.PID)
		}
	case "r", "R":
		m.statusMsg = ""
		return m, loadPortsCmd()
	}
	return m, nil
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		if m.cursor > 0 {
			m.cursor--
			m.adjustListOffset()
		}
	case tea.MouseButtonWheelDown:
		if m.cursor < len(m.ports)-1 {
			m.cursor++
			m.adjustListOffset()
		}
	case tea.MouseButtonLeft:
		if msg.Action != tea.MouseActionPress {
			break
		}
		startY := titleBarLines + 1 + listHeaderLines
		if msg.Y >= startY {
			if idx := (msg.Y - startY) + m.listOffset; idx < len(m.ports) {
				m.cursor = idx
				m.adjustListOffset()
			}
		}
	}
	return m, nil
}

func (m *Model) adjustListOffset() {
	h := m.listRowHeight()
	if h <= 0 {
		return
	}
	if m.cursor >= len(m.ports) {
		m.cursor = max(0, len(m.ports)-1)
	}
	if m.cursor < m.listOffset {
		m.listOffset = m.cursor
	} else if m.cursor >= m.listOffset+h {
		m.listOffset = m.cursor - h + 1
	}
}

func (m Model) listRowHeight() int {
	if h := m.height - titleBarLines - statusBarLines - panelBorderLines - listHeaderLines; h > 0 {
		return h
	}
	return 1
}

func tickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func loadPortsCmd() tea.Cmd {
	return func() tea.Msg {
		return portsMsg(network.ScanPorts())
	}
}

func killProcess(pid int) tea.Cmd {
	return func() tea.Msg {
		proc, err := os.FindProcess(pid)
		if err != nil {
			return killDoneMsg{pid: pid, err: err}
		}
		if err := proc.Signal(syscall.SIGTERM); err != nil {
			return killDoneMsg{pid: pid, err: err}
		}
		time.Sleep(500 * time.Millisecond)
		if err := proc.Signal(syscall.Signal(0)); err == nil {
			_ = proc.Signal(syscall.SIGKILL)
		}
		return killDoneMsg{pid: pid, err: nil}
	}
}

func selectedKey(ports []models.PortEntry, idx int) string {
	if idx < 0 || idx >= len(ports) {
		return ""
	}
	return portEntryKey(ports[idx])
}

func portEntryKey(p models.PortEntry) string {
	return strconv.Itoa(int(p.Port)) + "|" + p.Protocol + "|" + strconv.Itoa(p.PID)
}
