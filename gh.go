package ghdl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
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
}

func (gr *GHRelease) GetGHReleases() (string, string, error) {
	var tag string
	if gr.TagName == "" {
		tag = "latest"
	} else {
		tag = "tags/" + gr.TagName
	}
	binaryName := filepath.Base(gr.RepoPath)
	apiUrl := fmt.Sprint("https://api.github.com/repos/", gr.RepoPath, "/releases/", tag)

	// Get releases info
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	} else if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("requst to %v failed: %v", apiUrl, resp.Status)
	}
	defer resp.Body.Close()
	byte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	var respJSON APIReleaseResp
	if err := json.Unmarshal(byte, &respJSON); err != nil {
		return "", "", err
	}
	releaseAssets := respJSON.Assets
	if len(releaseAssets) == 0 {
		return "", "", fmt.Errorf("no binary release found")
	}

	// Pick release assets
	matchedAssets := filterAssets(filterAssets(releaseAssets, OS), ARCH)
	if len(matchedAssets) != 1 {
		var choices []string
		for _, asset := range matchedAssets {
			choices = append(choices, asset.Name)
		}
		idx := sl.Select(&choices)
		return matchedAssets[idx].DownloadUrl, binaryName, nil
	}
	return matchedAssets[0].DownloadUrl, binaryName, nil
}

// Filter assets by match pattern, falling back to the default assets if no match is found
func filterAssets(assets []APIReleaseAsset, match string) (ret []APIReleaseAsset) {
	for _, asset := range assets {
		if strings.Contains(strings.ToLower(asset.Name), match) {
			ret = append(ret, asset)
		}
	}
	if len(ret) == 0 {
		return assets
	}
	return ret
}
