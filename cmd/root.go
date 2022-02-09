package main

import (
	"os"
	"strings"

	"github.com/beetcb/ghdl"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-dl <user/repo[#tagname]>",
	Short: "gh-dl download binary from github release",
	Long: `gh-dl download binary from github release
gh-dl handles archived or compressed file as well`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, tag := parseArg(args[0])
		ghRelease := ghdl.GHRelease{RepoPath: repo, TagName: tag}
		url, binaryName, err := ghRelease.GetGHReleases()
		ghReleaseDl := ghdl.GHReleaseDl{Url: url, BinaryName: binaryName}
		if err != nil {
			panic(err)
		}
		binaryNameFlag, err := cmd.Flags().GetString("bin")
		if err != nil {
			panic(err)
		}
		if binaryNameFlag != "" {
			binaryName = binaryNameFlag
		}
		ghReleaseDl.DlAndDecompression()
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("bin", "b", "", "specify bin file name")
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
