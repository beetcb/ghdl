package ghdl

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/beetcb/ghdl/helper/pg"
)

type GHReleaseDl struct {
	BinaryName string
	Url        string
}

func (dl *GHReleaseDl) DlAndDecompression() {
	b := dl.BinaryName + func() string {
		if runtime.GOOS == "windows" {
			return ".exe"
		} else {
			return ""
		}
	}()

	req, err := http.NewRequest("GET", dl.Url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fileSize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		panic(err)
	}

	file, err := ioutil.TempFile("", "*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	// create progress tui
	starter := func(updater func(float64)) {
		if _, err := io.Copy(file, &pg.ProgressBytesReader{Reader: resp.Body, Handler: func(p int) {
			updater(float64(p) / float64(fileSize))
		}}); err != nil {
			panic(err)
		}
	}
	pg.Progress(starter)

	// `file` has no content, we must open it for reading
	openfile, err := os.Open(file.Name())
	if err != nil {
		panic(err)
	}
	filebody, err := ioutil.ReadAll(openfile)
	if err != nil {
		panic(err)
	}
	bytesReader := bytes.NewReader(filebody)

	fileExt := filepath.Ext(dl.Url)
	var decompressedBinary *[]byte
	switch fileExt {
	case ".zip":
		decompressedBinary, err = dl.zipBinary(bytesReader, b)
		if err != nil {
			panic(err)
		}
	case ".gz":
		if strings.Contains(dl.Url, ".tar.gz") {
			decompressedBinary, err = dl.targzBinary(bytesReader, b)
			if err != nil {
				panic(err)
			}
		} else {
			decompressedBinary, err = dl.gzBinary(bytesReader, b)
			if err != nil {
				panic(err)
			}
		}
	case ".deb":
	case ".rpm":
	case ".apk":
		fileName := b + fileExt
		fmt.Printf("Detected deb/rpm/apk package, download directly to ./%s\nYou can install it with the appropriate commands\n", fileName)
		if err := os.WriteFile(fileName, filebody, 0777); err != nil {
			panic(err)
		}
	case "":
		decompressedBinary = &filebody
	default:
		panic("unsupported file format")
	}
	if err := os.WriteFile(b, *decompressedBinary, 0777); err != nil {
		panic(err)
	}
}

func (*GHReleaseDl) zipBinary(r *bytes.Reader, b string) (*[]byte, error) {
	zipR, err := zip.NewReader(r, int64(r.Len()))
	if err != nil {
		return nil, err
	}

	for _, f := range zipR.File {
		if filepath.Base(f.Name) == b || len(zipR.File) == 1 {
			open, err := f.Open()
			if err != nil {
				return nil, err
			}
			ret, err := ioutil.ReadAll(open)
			if err != nil {
				return nil, err
			}
			return &ret, err
		}
	}
	return nil, fmt.Errorf("Binary file %v not found", b)
}

func (*GHReleaseDl) gzBinary(r *bytes.Reader, b string) (*[]byte, error) {
	gzR, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzR.Close()
	ret, err := ioutil.ReadAll(gzR)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (*GHReleaseDl) targzBinary(r *bytes.Reader, b string) (*[]byte, error) {
	gzR, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzR.Close()
	tarR := tar.NewReader(gzR)

	var file []byte
	for {
		header, err := tarR.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if (header.Typeflag != tar.TypeDir) && filepath.Base(header.Name) == b {
			file, err = ioutil.ReadAll(tarR)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	return &file, nil
}
