package ghdl

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/beetcb/ghdl/helper/pg"
	humanize "github.com/dustin/go-humanize"
)

var (
	ErrNeedInstall = errors.New(
		"detected deb/rpm/apk package, download directly")
	ErrNoBin = errors.New("binary file not found")
)

type GHReleaseDl struct {
	BinaryName string
	Url        string
	Size       int64
}

// Download asset from github release to `path`
//
// dl.BinaryName shall be replaced with absolute path mutably
func (dl *GHReleaseDl) DlTo(path string) (err error) {
	dl.BinaryName, err = filepath.Abs(filepath.Join(path, dl.BinaryName))
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", dl.Url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpfile, err := os.Create(dl.BinaryName + ".tmp")
	if err != nil {
		return err
	}
	defer tmpfile.Close()

	// create progress tui
	starter := func(updater func(float64)) {
		if _, err := io.Copy(tmpfile, &pg.ProgressBytesReader{Reader: resp.Body, Handler: func(p int) {
			updater(float64(p) / float64(dl.Size))
		}}); err != nil {
			panic(err)
		}
	}
	pg.Progress(starter, humanize.Bytes(uint64(dl.Size)))
	return nil
}

// Extract binary file from the downloaded temporary file.
//
// Currently supporting unarchiving `tar` and decompressing `zip` `gravezip`.
//
// Package format `deb` `rpm` `apk` will be downloaded directly
func (dl GHReleaseDl) ExtractBinary() error {
	tmpfileName := dl.BinaryName + ".tmp"
	openfile, err := os.Open(tmpfileName)
	if err != nil {
		return err
	}

	fileExt := filepath.Ext(dl.Url)
	var decompressedBinary io.Reader
	switch fileExt {
	case ".zip":
		zipFile, err := dl.UnZipBinary(openfile)
		if err != nil {
			return err
		}
		decompressedBinary, err = zipFile.Open()
		if err != nil {
			return err
		}
	case ".gz":
		if strings.Contains(dl.Url, ".tar.gz") {
			decompressedBinary, err = dl.UnTargzBinary(openfile)
			if err != nil {
				return err
			}
		} else {
			decompressedBinary, err = dl.UnGzBinary(openfile)
			if err != nil {
				return err
			}
		}
	case "":
		decompressedBinary = openfile
	case ".deb", ".rpm", ".apk", ".msi", ".exe", ".dmg":
		fileName := dl.BinaryName + fileExt
		if err := os.Rename(tmpfileName, fileName); err != nil {
			panic(err)
		}
		return ErrNeedInstall
	default:
		defer os.Remove(tmpfileName)
		return fmt.Errorf("unsupported file format: %v", fileExt)
	}
	defer os.Remove(tmpfileName)
	defer openfile.Close()
	out, err := os.Create(dl.BinaryName)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, decompressedBinary); err != nil {
		return err
	}
	return nil
}

func (dl GHReleaseDl) UnZipBinary(r *os.File) (*zip.File, error) {
	b := filepath.Base(dl.BinaryName)
	zipR, err := zip.NewReader(r, dl.Size)
	if err != nil {
		return nil, err
	}

	for _, f := range zipR.File {
		if filepath.Base(f.Name) == b || len(zipR.File) == 1 {
			return f, nil
		}
	}
	return nil, ErrNoBin
}

func (GHReleaseDl) UnGzBinary(r *os.File) (*gzip.Reader, error) {
	gzR, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzR.Close()
	return gzR, nil
}

func (dl GHReleaseDl) UnTargzBinary(r *os.File) (*tar.Reader, error) {
	b := filepath.Base(dl.BinaryName)
	gzR, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzR.Close()
	tarR := tar.NewReader(gzR)

	for {
		header, err := tarR.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if (header.Typeflag != tar.TypeDir) && filepath.Base(header.Name) == b {
			if err != nil {
				return nil, err
			}
			return tarR, nil
		}
	}
	return nil, ErrNoBin
}
