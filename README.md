# ghdl

> Memorize `ghdl` as `github download`

`ghdl` is a fast and simple program (and also a golang module) for downloading and installing executable binary from github releases.

<p align="center">
    <img alt="animated demo" src="./demo.svg" width="600px">
</p>
<p align="center">
  <strong>The demo above extracts <code>fd</code> execuable to current working directory and give execute permission to it.</strong>
</p>

# Features

- Auto decompressing and unarchiving the downloaded asset (without any system dependencies like `tar` or `unzip`).

    ```ts
    Currently supporting unarchiving `tar` and decompressing `zip` `gzip`.
    Package format `deb` `rpm` `apk` will be downloaded directly
    ```
- Setups for executable: `ghdl` moves executable to specified location and add execute permissions to the file.
- Auto filtering: multiple assets in one release will be filtered by OS or ARCH.
- Interactive TUI: when auto filtering is failed or returned multiple options, you can select assets in a interactive way, with vim key bindings support.
- Release tags: `ghdl` downloads latest release by default, other or old tagged releases can be downloaded by specifying release tag: `username/repo#tagname`
- See download with real-time progress bar.

# Installation

> If you're going to use `ghdl` as a go module, ignore the following installation progress.

- Using Go tools: 

    go will download the latest version of ghdl to $GOPATH/bin, please make sure $GOPATH is in the PATH: 

    ```sh
    go install github.com/beetcb/ghdl/ghdl@latest
    ```

- Download and run executable from release.
- Run the following shell script(*nix system only):

    ```sh
    curl -fsSL "https://bina.egoist.sh/beetcb/ghdl?dir=/usr/local/bin" | sh
    # feel free to change the `dir` url param to specify the installation directory.
    ```

# Usage

### CLI
Run `ghdl --help`

```sh
ghdl --help
    ghdl download binary from github release
    ghdl handles archived or compressed file as well

    Usage:
    ghdl <user/repo[#tagname]> [flags]

    Flags:
    -h, --help          help for ghdl
    -n, --name string   specify binary file name
    -p, --path path     save binary to path (default ".")

```

It's tedious to specify `-p` manually, we can alias `ghdl -p "$DirInPath"` to a shorthand command, then use it as a executable installer.

### Go Module

Require `ghdl` to go.mod

```sh
go get github.com/beetcb/ghdl
```

Use `ghdl`'s out-of-box utilities:

```go
package main

import (
	"fmt"
	"os"

	"github.com/beetcb/ghdl"
	h "github.com/beetcb/ghdl/helper"
)

func main() {
	ghRelease := ghdl.GHRelease{RepoPath: "sharkdp/fd"}
	ghReleaseDl, err := ghRelease.GetGHReleases()
	if err != nil {
		h.Print(fmt.Sprintf("get gh releases failed: %s", err), h.PrintModeErr)
		os.Exit(1)
	}

	if err := ghReleaseDl.DlTo("."); err != nil {
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
}
```

# Credit

Inspired by [egoist/bina](https://github.com/egoist/bina), TUI powered by [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)

# License

Licensed under [MIT](./LICENSE)

Author: @beetcb | Email: i@beetcb.com
