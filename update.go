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
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/mod/semver"
)

const (
	urlDbUpdate = "https://www.wireshark.org/download/automated/data/manuf.gz"
	githubOwner = "adhoniran"
	githubRepo  = "ouimap"
	urlBaseRepo = "https://github.com/%s/%s"
	urlBaseAPI  = "https://api.github.com/repos/%s/%s/releases/latest"
)

var (
	exePath, _ = os.Executable()
	dbPath     = filepath.Join(filepath.Dir(exePath), "manuf.gz")
	dbTemp     = filepath.Join(filepath.Dir(exePath), "manuf.tmp")
)

func downloadDatabase(url, tmp string) error {
	resp, err := http.Get(url)
	if err != nil {
		// return fmt.Errorf("HTTP request failed: %w", err)
		return fmt.Errorf("%s!%s unable to check for database updates", red, reset)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close HTTP response: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	tmpFile, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		if closeErr := tmpFile.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close temp file: %v\n", closeErr)
		}
	}()

	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetDescription("Downloading OUI database"),
		progressbar.OptionShowCount(),
		progressbar.OptionShowBytes(true),
	)

	if _, err = io.Copy(io.MultiWriter(tmpFile, bar), resp.Body); err != nil {
		return fmt.Errorf("failed to save data to temp file: %w", err)
	}

	return nil
}

func checkDatabase(file string) error {
	db, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	gzipReader, err := gzip.NewReader(db)
	if err != nil {
		return fmt.Errorf("failed to initialize gzip reader: %w", err)
	}
	defer func() {
		if cerr := gzipReader.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close gzip reader: %w", cerr)
		}
	}()

	if _, err = io.Copy(io.Discard, gzipReader); err != nil {
		return fmt.Errorf("database integrity check failed: %w", err)
	}

	return nil
}

func replaceDatabase(db, tmp string) error {
	if _, err := os.Stat(db); err == nil {
		if err := os.Remove(db); err != nil {
			return fmt.Errorf("failed to remove old database: %w", err)
		}
	}
	if err := os.Rename(tmp, db); err != nil {
		return fmt.Errorf("failed to rename temp file to new database: %w", err)
	}
	return nil
}

func updateDatabase() (bool, error) {

	info, err := os.Stat(dbPath)
	if err == nil {
		if time.Since(info.ModTime()) < 7*24*time.Hour {
			return false, nil
		}
	} else if !os.IsNotExist(err) {
		return false, err
	}

	if err := downloadDatabase(urlDbUpdate, dbTemp); err != nil {
		return false, fmt.Errorf("failed to download OUI database: %w", err)
	}

	if err := checkDatabase(dbTemp); err != nil {
		return false, fmt.Errorf("integrity check failed: %w", err)
	}

	if err := replaceDatabase(dbPath, dbTemp); err != nil {
		return false, fmt.Errorf("failed to replace database: %w", err)
	}

	return true, nil
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func checkNewVersion(owner, repo, currentVersion string) (bool, string) {
	url := fmt.Sprintf(urlBaseAPI, owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// fmt.Printf("HTTP request failed: %v\n", err)
		fmt.Printf("%s!%s unable to check for software updates\n", red, reset)
		return false, ""
	}
	httpUserAgent := fmt.Sprintf("OUImap/%s", version)
	req.Header.Set("User-Agent", httpUserAgent)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		// fmt.Printf("HTTP request failed: %v\n", err)
		fmt.Printf("%s!%s unable to check for database updates\n", red, reset)
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
		url := fmt.Sprintf(urlBaseRepo, githubOwner, githubRepo)
		fmt.Printf("Check for new version at %s\n", url)
		return false, latestVersion

	}

	return false, latestVersion
}
