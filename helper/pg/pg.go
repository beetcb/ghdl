package pg

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type ProgressBytesReader struct {
	io.Reader
	Progressed int
	Handler    func(p int)
}

type model struct {
	percent  float64
	humanize string
	progress progress.Model
	init     func() tea.Msg
}

func (pbr *ProgressBytesReader) Read(b []byte) (n int, err error) {
	if n, err = pbr.Reader.Read(b); err != nil {
		return n, err
	}
	pbr.Progressed += n
	pbr.Handler(pbr.Progressed)
	return
}

const (
	padding  = 2
	maxWidth = 80
)

func (m model) Init() tea.Cmd {
	return m.init
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.percent == 1 {
		return m, tea.Quit
	}
	return m, func() tea.Msg { m.Update(nil); return nil }
}

func (e model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + e.progress.ViewAs(e.percent) + fmt.Sprintf(" of %s", e.humanize) + "\n"
}

func Progress(starter func(updater func(float64)), humanize string) {
	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	state := model{progress: prog, humanize: humanize}
	updater := func(p float64) {
		state.percent = p
	}
	state.init = func() tea.Msg {
		starter(updater)
		return nil
	}

	if err := tea.NewProgram(&state).Start(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}
