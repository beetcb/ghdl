package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
		filterOff, err := cmdFlags.GetBool("filter-off")
		if err != nil {
			panic(err)
		}
		assetFilterString, err := cmdFlags.GetString("asset-filter")
		if err != nil {
			panic(err)
		}
		var assetFilter *regexp.Regexp
		if assetFilterString != "" {
			assetFilter, err = regexp.Compile(assetFilterString)
			if err != nil {
				panic(err)
			}
		}

		repo, tag := parseArg(args[0])
		ghRelease := ghdl.GHRelease{RepoPath: repo, TagName: tag}
		ghReleaseDl, err := ghRelease.GetGHReleases(filterOff, assetFilter)

		if err != nil {
			h.Println(fmt.Sprintf("get gh releases failed: %s", err), h.PrintModeErr)
			os.Exit(1)
		}

		if binaryNameFlag != "" {
			ghReleaseDl.BinaryName = binaryNameFlag
		}
		h.Println(fmt.Sprintf("start downloading %s", h.Sprint(filepath.Base(ghReleaseDl.Url), h.SprintOptions{PromptOff: true, PrintMode: h.PrintModeSuccess})), h.PrintModeInfo)
		if err := ghReleaseDl.DlTo(pathFlag); err != nil {
			h.Println(fmt.Sprintf("download failed: %s", err), h.PrintModeErr)
			os.Exit(1)
		}
		if err := ghReleaseDl.ExtractBinary(); err != nil {
			switch err {
			case ghdl.ErrNeedInstall:
				h.Println(fmt.Sprintf("%s. You can install it with the appropriate commands", err), h.PrintModeInfo)
				os.Exit(0)
			case ghdl.ErrNoBin:
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
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("name", "n", "", "specify binary file name to enhance filtering and extracting accuracy")
	rootCmd.PersistentFlags().StringP("asset-filter", "f", "",
		"specify regular expression for the asset name; used in conjunction with the platform and architecture filters.")
	rootCmd.PersistentFlags().StringP("path", "p", ".", "save binary to `path` and add execute permission to it")
	rootCmd.PersistentFlags().BoolP("filter-off", "F", false, "turn off auto-filtering feature")
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
