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
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	dbUpdateURL    = "https://www.wireshark.org/download/automated/data/manuf.gz"
	dbDownloadFile = "ouimap.dl"
	exePath, _     = os.Executable()
	dbStorageFile  = filepath.Join(filepath.Dir(exePath), "ouimap.db")
)

func downloadDatabase(url, tmp string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
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

	if _, err = io.Copy(tmpFile, resp.Body); err != nil {
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

	update := false

	info, err := os.Stat(dbStorageFile)
	if err == nil {
		if time.Since(info.ModTime()) < 7*24*time.Hour {
			return update, nil
		}
	} else if !os.IsNotExist(err) {

		return update, err
	}

	if err := downloadDatabase(dbUpdateURL, dbDownloadFile); err != nil {
		return update, fmt.Errorf("failed to download OUI database: %w", err)
	}

	if err := checkDatabase(dbDownloadFile); err != nil {
		return update, fmt.Errorf("integrity check failed: %w", err)
	}

	if err := replaceDatabase(dbStorageFile, dbDownloadFile); err != nil {
		return update, fmt.Errorf("failed to replace database: %w", err)
	}

	update = true
	return update, nil
}
