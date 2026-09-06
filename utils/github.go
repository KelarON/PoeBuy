package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// GetVersion returns the version of the release
func (r *ReleaseInfo) GetVersion() string {
	return r.TagName
}

// GetLatestRelease fetches the latest release information from the GitHub API for the PoeBuy repository
func GetLatestRelease() (*ReleaseInfo, error) {

	// Prepare the HTTP request to GitHub API
	request, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/KelarON/PoeBuy/releases/latest", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Add("X-GitHub-Api-Version", "2026-03-10")

	// Send the request and get the response
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch latest release: %s", response.Status)
	}

	// Parse the JSON response
	releaseInfo := new(ReleaseInfo)
	err = json.NewDecoder(response.Body).Decode(releaseInfo)
	if err != nil {
		return nil, err
	}

	return releaseInfo, nil
}

// download fetches the PoeBuy.exe asset from the release and returns its content as a byte slice
func (r *ReleaseInfo) download() ([]byte, error) {

	// Find the asset with the PoeBuy.exe name
	var downloadURL string

	for _, asset := range r.Assets {
		if asset.Name == UPDATE_FILE_NAME {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return nil, fmt.Errorf("no PoeBuy.exe asset found")
	}

	// Download the release
	response, err := http.Get(downloadURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download release: %s", response.Status)
	}

	return io.ReadAll(response.Body)
}

// SaveRelease saves the downloaded release data to the specified file path
func SaveRelease(release *ReleaseInfo, filePath string) error {

	data, err := release.download()
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
