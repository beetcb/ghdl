package ghdl

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/beetcb/ghdl/helper/pg"
	humanize "github.com/dustin/go-humanize"
)

type GHReleaseDl struct {
	BinaryName string
	Url        string
	Size       int64
}

// Download asset from github release
// dl.BinaryName path might change mutably
func (dl *GHReleaseDl) DlTo(path string) error {
	if path != "" {
		dl.BinaryName = filepath.Join(path, dl.BinaryName)
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

	if err != nil {
		return err
	}

	file, err := os.Create(dl.BinaryName)
	if err != nil {
		return err
	}

	// create progress tui
	starter := func(updater func(float64)) {
		if _, err := io.Copy(file, &pg.ProgressBytesReader{Reader: resp.Body, Handler: func(p int) {
			updater(float64(p) / float64(dl.Size))
		}}); err != nil {
			panic(err)
		}
	}
	pg.Progress(starter, humanize.Bytes(uint64(dl.Size)))
	return nil
}

func (dl GHReleaseDl) ExtractBinary() error {
	// `file` has no content, we must open it for reading
	openfile, err := os.Open(dl.BinaryName)
	if err != nil {
		return err
	}
	defer openfile.Close()

	fileExt := filepath.Ext(dl.Url)
	var decompressedBinary io.Reader = nil
	switch fileExt {
	case ".zip":
		decompressedBinary, err = dl.ZipBinary(openfile)
		if err != nil {
			return err
		}
	case ".gz":
		if strings.Contains(dl.Url, ".tar.gz") {
			decompressedBinary, err = dl.TargzBinary(openfile)
			if err != nil {
				return err
			}
		} else {
			decompressedBinary, err = dl.GzBinary(openfile)
			if err != nil {
				return err
			}
		}
	case "":
		decompressedBinary = openfile
	case ".deb":
	case ".rpm":
	case ".apk":
		fileName := dl.BinaryName + fileExt
		fmt.Printf("Detected deb/rpm/apk package, download directly to ./%s\nYou can install it with the appropriate commands\n", fileName)
		if err := os.Rename(dl.BinaryName, fileName); err != nil {
			panic(err)
		}
		return nil
	default:
		defer os.Remove(dl.BinaryName)
		return fmt.Errorf("unsupported file format")
	}

	// rewrite the file
	out, err := os.Create(dl.BinaryName)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, decompressedBinary); err != nil {
		return nil
	}

	return nil
}

func (dl GHReleaseDl) ZipBinary(r *os.File) (io.Reader, error) {
	zipR, err := zip.NewReader(r, dl.Size)
	if err != nil {
		return nil, err
	}

	for _, f := range zipR.File {
		if filepath.Base(f.Name) == dl.BinaryName || len(zipR.File) == 1 {
			open, err := f.Open()
			if err != nil {
				return nil, err
			}
			return open, nil
		}
	}
	return nil, fmt.Errorf("Binary file %v not found", dl.BinaryName)
}

func (GHReleaseDl) GzBinary(r *os.File) (io.Reader, error) {
	gzR, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzR.Close()
	return gzR, nil
}

func (dl GHReleaseDl) TargzBinary(r *os.File) (io.Reader, error) {
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
		if (header.Typeflag != tar.TypeDir) && filepath.Base(header.Name) == dl.BinaryName {

			if err != nil {
				return nil, err
			}
			break
		}
	}
	return tarR, nil
}
