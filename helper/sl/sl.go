package sl

import (
	"fmt"
	"os"

	h "github.com/beetcb/ghdl/helper"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 80

type model struct {
	choices  []string // items on the to-do list
	cursor   int      // which to-do list item our cursor is pointing at
	selected int      // which to-do items are selected
}

func initialModel(choices *[]string) model {
	return model{
		choices:  *choices,
		selected: -1,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.cursor
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	blue, printWidth := lipgloss.Color("14"), 60
	paddingS := lipgloss.NewStyle().PaddingLeft(2).Width(printWidth)
	colorS := paddingS.Copy().
		Foreground(blue).BorderLeft(true).BorderForeground(blue)
	s := "\n" + h.Sprint("there is more than one option after filtering, please select it manually", h.PrintModeInfo) + "\n"
	if m.selected == -1 {
		s += "\n"
		for i, choice := range m.choices {
			if m.cursor == i {
				s += colorS.Render(choice) + "\n"
			} else {
				s += paddingS.Render(choice) + "\n"
			}
		}
		// Send the UI for rendering
		return s + "\n"
	} else {
		s += h.Sprint(fmt.Sprintf("start downloading %s", lipgloss.NewStyle().Foreground(blue).Render(m.choices[m.selected])), h.PrintModeInfo) + "\n"
		return s
	}
}

func Select(choices *[]string) int {
	state := initialModel(choices)
	p := tea.NewProgram(&state)
	if err := p.Start(); err != nil {
		h.Print(fmt.Sprintf("Alas, there's been an error: %v", err), h.PrintModeErr)
		os.Exit(1)
	}
	return state.selected
}
