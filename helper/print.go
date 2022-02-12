package helper

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	PrintModeInfo    = 0
	PrintModeSuccess = 1
	PrintModeErr     = 2
)

type SprintOptions struct {
	PrintMode  int
	PrintWidth int
	PromptOff  bool
}

func Sprint(str string, options SprintOptions) string {
	newStyle := lipgloss.NewStyle()
	if options.PrintWidth > 0 {
		newStyle = newStyle.Width(options.PrintWidth)
	}
	sPrint := ""
	if !options.PromptOff {
		prompt := lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Render("→ ")
		sPrint = prompt
	}

	switch options.PrintMode {
	case PrintModeInfo:
		sPrint += newStyle.Copy().Foreground(lipgloss.Color("146")).Render(str)
	case PrintModeSuccess:
		sPrint += newStyle.Copy().Foreground(lipgloss.Color("6")).Render(str)
	case PrintModeErr:
		sPrint += newStyle.Copy().Foreground(lipgloss.Color("9")).Render(str)
	}

	return sPrint
}

func Println(str string, printMode int) {
	sPrint := Sprint(str, SprintOptions{PrintMode: printMode})
	fmt.Println(sPrint)
}
