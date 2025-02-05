package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

var (
	githubOwner = "adhoniran"
	githubRepo  = "ouimap"
)

func checkNewVersion(owner, repo, currentVersion string) (bool, string) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, ""
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, ""
	}
	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, ""
	}
	latestVersion := release.TagName
	isNewer := strings.TrimPrefix(latestVersion, "v") != strings.TrimPrefix(currentVersion, "v")
	return isNewer, latestVersion
}
