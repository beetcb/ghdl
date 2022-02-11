# ghdl

`ghdl`(short for github download) is a fast and simple program for downloading and installing executable binary from GitHub releases.

# Features

![demo](./demo.svg)

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
- Using Go tools: 

    go will download the latest version of ghdl to $GOPATH/bin, please make sure $GOPATH is in the PATH: 

    ```sh
    go install github.com/beetcb/ghdl/ghdl@latest`
    ```

- Download and run executable from release.
- Run the following shell script(*nix system only):

    ```sh
    curl -fsSL "https://bina.egoist.sh/beetcb/ghdl?dir=/usr/local/bin" | sh
    # feel free to change the `dir` url param to specify the installation directory.
    ```

# Usage

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

# Credit

Inspired by [egoist/bina](https://github.com/egoist/bina), TUI powered by [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)

# License

Licensed under [MIT](./LICENSE)

Author: @beetcb | Email: i@beetcb.com
