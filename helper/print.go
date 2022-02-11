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

func Sprint(str string, printMode int) string {
	printWidth := 60
	var newStyle = lipgloss.NewStyle().Width(printWidth)
	prompt := lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Render("→ ")
	var sPrint string = prompt
	switch printMode {
	case PrintModeInfo:
		sPrint += newStyle.Copy().Foreground(lipgloss.Color("146")).Render(str)
	case PrintModeSuccess:
		sPrint += newStyle.Copy().Foreground(lipgloss.Color("6")).Render(str)
	case PrintModeErr:
		sPrint += newStyle.Copy().Foreground(lipgloss.Color("9")).Render(str)
	}

	return sPrint
}

func Print(str string, printMode int) {
	sPrint := Sprint(str, printMode)
	fmt.Println(sPrint)
}
