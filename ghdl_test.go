package ghdl

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	h "github.com/beetcb/ghdl/helper"
)

// Shall never failed (if complied) but can test TUI behavior and functionality
func TestDownloadFdBinary(t *testing.T) {
	ghRelease := GHRelease{RepoPath: "sharkdp/fd"}
	ghReleaseDl, err := ghRelease.GetGHReleases(false)

	if err != nil {
		h.Println(fmt.Sprintf("get gh releases failed: %s", err), h.PrintModeErr)
		os.Exit(1)
	}

	h.Println(fmt.Sprintf("start downloading %s", h.Sprint(filepath.Base(ghReleaseDl.Url), h.SprintOptions{PromptOff: true, PrintMode: h.PrintModeSuccess})), h.PrintModeInfo)
	if err := ghReleaseDl.DlTo("."); err != nil {
		h.Println(fmt.Sprintf("download failed: %s", err), h.PrintModeErr)
		os.Exit(1)
	}
	if err := ghReleaseDl.ExtractBinary(); err != nil {
		switch err {
		case ErrNeedInstall:
			h.Println(fmt.Sprintf("%s. You can install it with the appropriate commands", err), h.PrintModeInfo)
			os.Exit(0)
		case ErrNoBin:
			h.Println(fmt.Sprintf("%s. Try to specify binary name flag", err), h.PrintModeInfo)
			os.Exit(0)
		default:
			h.Println(fmt.Sprintf("extract failed: %s", err), h.PrintModeErr)
			os.Exit(1)
		}
	}
	h.Println(fmt.Sprintf("saved executable to %s", ghReleaseDl.BinaryName), h.PrintModeSuccess)
	if err := os.Chmod(ghReleaseDl.BinaryName, 0777); err != nil {
		h.Println(fmt.Sprintf("chmod failed: %s", err), h.PrintModeErr)
	}
}
