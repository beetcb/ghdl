package helper

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 80

const (
	PrintModeInfo    = 0
	PrintModeSuccess = 1
	PrintModeErr     = 2
)

func Print(str string, printMode int) {
	var PaddingLeft = lipgloss.NewStyle().PaddingLeft(2).MaxWidth(maxWidth)
	switch printMode {
	case PrintModeInfo:
		fmt.Println(PaddingLeft.Foreground(lipgloss.Color("11")).Render(str))
	case PrintModeSuccess:
		fmt.Println(PaddingLeft.Foreground(lipgloss.Color("14")).Render(str))
	case PrintModeErr:
		fmt.Println(PaddingLeft.Foreground(lipgloss.Color("202")).Render(str))
	}
}
