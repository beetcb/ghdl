package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/beetcb/ghdl"
	h "github.com/beetcb/ghdl/helper"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghdl <user/repo[#tagname]>",
	Short: "ghdl download binary from github release",
	Long: `ghdl download binary from github release
ghdl handles archived or compressed file as well`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdFlags := cmd.Flags()
		binaryNameFlag, err := cmdFlags.GetString("name")
		if err != nil {
			panic(err)
		}
		pathFlag, err := cmdFlags.GetString("path")
		if err != nil {
			panic(err)
		}
		repo, tag := parseArg(args[0])
		ghRelease := ghdl.GHRelease{RepoPath: repo, TagName: tag}
		ghReleaseDl, err := ghRelease.GetGHReleases()
		if err != nil {
			h.Print(fmt.Sprintf("get gh releases failed: %s", err), h.PrintModeErr)
			os.Exit(1)
		}

		if binaryNameFlag != "" {
			ghReleaseDl.BinaryName = binaryNameFlag
		}
		if err := ghReleaseDl.DlTo(pathFlag); err != nil {
			h.Print(fmt.Sprintf("download failed: %s", err), h.PrintModeErr)
			os.Exit(1)
		}
		if err := ghReleaseDl.ExtractBinary(); err != nil {
			switch err {
			case ghdl.NeedInstallError:
				h.Print(fmt.Sprintf("%s. You can install %s with the appropriate commands", err, ghReleaseDl.BinaryName), h.PrintModeInfo)
				os.Exit(0)
			case ghdl.NoBinError:
				h.Print(fmt.Sprintf("%s. Try to specify binary name flag", err), h.PrintModeInfo)
				os.Exit(0)
			default:
				h.Print(fmt.Sprintf("extract failed: %s", err), h.PrintModeErr)
				os.Exit(1)
			}
		}
		h.Print(fmt.Sprintf("saved executable to %s", ghReleaseDl.BinaryName), h.PrintModeSuccess)
		if err := os.Chmod(ghReleaseDl.BinaryName, 0777); err != nil {
			h.Print(fmt.Sprintf("chmod failed: %s", err), h.PrintModeErr)
		}
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("name", "n", "", "specify binary file name")
	rootCmd.PersistentFlags().StringP("path", "p", ".", "save binary to `path`")
}

// parse user/repo[#tagname] arg
func parseArg(repoPath string) (repo string, tag string) {
	seperateTag := strings.Split(repoPath, "#")
	if len(seperateTag) == 2 {
		tag = seperateTag[1]
	}
	repo = seperateTag[0]
	return repo, tag
}
