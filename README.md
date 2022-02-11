# ghdl

`ghdl`(short for github download) is a fast and simple program for downloading and installing executable binary from GitHub releases.

# Features

- Auto decompressing and unarchiving the downloaded asset.

    ```ts
    Currently supporting unarchiving `tar` and decompressing `zip` `gzip`.

    Package format `deb` `rpm` `apk` will be downloaded directly
    ```
- Setups for executable: `ghdl` moves executable to specified location and add execute permissions to the file.
- Auto filtering: multiple assets in one release will be filtered by OS or ARCH.
- Interactive TUI: when auto filtering is failed or returned multiple options, you can select assets in a interactive way, with vim key bindings support.
- Release tags: `ghdl` downloads latest release by default, other or old tagged releases can be downloaded by specifying release tag: `username/repo#tagname`

# Installation
- Using Go tools: 

    go will download the latest version of ghdl to $GOPATH/bin, please make sure $GOPATH is in the PATH: 

    ```sh
    go install github.com/beetcb/ghdl/ghdl@latest`
    ```
- Download executable from release.


