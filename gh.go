package ghdl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/beetcb/ghdl/helper/sl"
)

const (
	OS   = runtime.GOOS
	ARCH = runtime.GOARCH
)

type GHRelease struct {
	RepoPath string
	TagName  string
}

type APIReleaseResp struct {
	Assets []APIReleaseAsset `json:"assets"`
}

type APIReleaseAsset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
	Size        int    `json:"size"`
}

type assetNamePredicate func(assetName string) bool

func (gr GHRelease) GetGHReleases(filterOff bool, assetFilter *regexp.Regexp) (*GHReleaseDl, error) {
	var tag string
	if gr.TagName == "" {
		tag = "latest"
	} else {
		tag = "tags/" + gr.TagName
	}

	// Os-specific binaryName
	binaryName := filepath.Base(gr.RepoPath) + func() string {
		if runtime.GOOS == "windows" {
			return ".exe"
		} else {
			return ""
		}
	}()
	apiUrl := fmt.Sprint("https://api.github.com/repos/", gr.RepoPath, "/releases/", tag)

	// Get releases info
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requst to %v failed: %v", apiUrl, resp.Status)
	}
	defer resp.Body.Close()
	byte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var respJSON APIReleaseResp
	if err := json.Unmarshal(byte, &respJSON); err != nil {
		return nil, err
	}
	releaseAssets := respJSON.Assets
	if len(releaseAssets) == 0 {
		return nil, fmt.Errorf("no binary release found")
	}

	// Filter or Pick release assets
	matchedAssets := func() []APIReleaseAsset {
		if filterOff {
			return releaseAssets
		}

		osArchPredicate := func(match string) assetNamePredicate {
			return func(assetName string) bool {
				return strings.Contains(assetName, match) ||
					(match == "amd64" && (strings.Contains(assetName, "x64") || strings.Contains(assetName, "x86_64")))
			}
		}
		predicates := []assetNamePredicate{
			osArchPredicate(OS),
			osArchPredicate(ARCH),
		}
		if assetFilter != nil {
			predicates = append(predicates, func(assetName string) bool {
				return assetFilter.MatchString(assetName)
			})
		}
		return filterAssets(releaseAssets, predicates)
	}()
	matchedIdx := 0
	if len(matchedAssets) != 1 {
		var choices []string
		for _, asset := range matchedAssets {
			choices = append(choices, asset.Name)
		}
		idx := sl.Select(&choices)
		matchedIdx = idx
	}
	asset := matchedAssets[matchedIdx]
	return &GHReleaseDl{binaryName, asset.DownloadUrl, int64(asset.Size)}, nil
}

// filterAssets assets using the provided predicates, falling back to the default assets if no match is found
func filterAssets(assets []APIReleaseAsset, predicates []assetNamePredicate) []APIReleaseAsset {
	ret := assets
	for _, p := range predicates {
		ret = filter(ret, p)
	}
	return ret
}

func filter(assets []APIReleaseAsset, predicate assetNamePredicate) []APIReleaseAsset {
	var ret []APIReleaseAsset
	for _, asset := range assets {
		lowerName := strings.ToLower(asset.Name)
		if predicate(lowerName) {
			ret = append(ret, asset)
		}
	}
	if len(ret) == 0 {
		return assets
	}
	return ret
}
