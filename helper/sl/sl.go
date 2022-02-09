package sl

import (
	"fmt"
	"os"

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
	blue := lipgloss.Color("14")
	paddingS := lipgloss.NewStyle().PaddingLeft(2).MaxWidth(maxWidth)
	colorS := paddingS.Copy().
		Foreground(blue).BorderLeft(true).BorderForeground(blue)
	if m.selected == -1 {
		s := "\n" + paddingS.Render("gh-dl can't figure out which release to download\nplease select it manully") + "\n\n"
		for i, choice := range m.choices {
			if m.cursor == i {
				s += colorS.Render(choice) + "\n"
			} else {
				s += paddingS.Render(choice) + "\n"
			}
		}
		// Send the UI for rendering
		return s
	} else {
		s := paddingS.Render(fmt.Sprintf("start downloading %s", lipgloss.NewStyle().Foreground(blue).Render(m.choices[m.selected])))
		return "\n" + s + "\n"
	}
}

func Select(choices *[]string) int {
	state := initialModel(choices)
	p := tea.NewProgram(&state)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	return state.selected
}
