// SPDX-FileCopyrightText: 2025 Adhoniran Gomes
// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileNotice:
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3 of the License only.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/mod/semver"
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
		fmt.Printf("HTTP request failed: %v\n", err)
		return false, ""
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return false, ""
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return false, ""
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		fmt.Printf("Failed to decode response: %v\n", err)
		return false, ""
	}

	latestVersion := release.TagName

	if !strings.HasPrefix(latestVersion, "v") {
		latestVersion = "v" + latestVersion
	}
	if !strings.HasPrefix(currentVersion, "v") {
		currentVersion = "v" + currentVersion
	}

	if semver.IsValid(latestVersion) && semver.IsValid(currentVersion) {
		cmp := semver.Compare(latestVersion, currentVersion)
		if cmp > 0 {
			return true, latestVersion
		}
	} else {
		fmt.Printf("Invalid version format. currentVersion: %s, latestVersion: %s\n", currentVersion, latestVersion)
	}

	return false, latestVersion
}
