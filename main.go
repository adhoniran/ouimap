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
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	appTitle     string = "OUImap"
	appArch             = runtime.GOARCH
	appOS               = runtime.GOOS
	appCopyright string = "Copyright © 2025 Adhoniran Gomes"
	appLicense   string = `This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it under certain conditions;
Visit https://www.gnu.org/licenses/gpl-3.0.html for details.`
)

var (
	appVersion   string
	buildVersion string
)

var (
	reset  string
	red    string
	green  string
	blue   string
	yellow string
)

func main() {

	fmt.Println()
	fmt.Println(blue + appTitle + " " + appVersion + "+build.g" + buildVersion + " (" + appOS + "/" + appArch + ") " + reset)
	fmt.Println(blue + appCopyright + reset)
	fmt.Println(blue + appLicense + reset)
	fmt.Println()

	updated, err := updateDatabase()
	if err != nil {
		fmt.Println(err)
	}
	if updated {
		fmt.Println()
		fmt.Println(green + "OUI database updated successfully!" + reset)
	}

	isNewer, latestVersion := checkNewVersion(githubOwner, githubRepo, appVersion)
	if isNewer {
		url := fmt.Sprintf(baseRepoURL, githubOwner, githubRepo)
		fmt.Printf("%sA new version of OUImap (%s) is available at %s\n%s", yellow, latestVersion, url, reset)
	}

	if err := loadOUIDatabase(dbPath); err != nil {
		fmt.Printf(red+"Failed to load OUI data: %s\n"+reset, err)
		return
	}

	if len(os.Args) == 1 {
		fmt.Print(`
Enter a multi-line list of OUIs, MAC addresses and/or descriptions. Separate OUI/MAC address parts with colons, hyphens or periods.
Press ENTER on a blank line to start the search, or CTRL+C to exit.


`)
	}

	promptContinue()

}

func promptContinue() {

	for {

		searchParams := getSearchParams()

		startTime := time.Now()
		searchParams = deduplicateInput(searchParams)

		lookupResults := lookupOUIData(searchParams)
		lookupResults = deduplicateResults(lookupResults)
		elapsedTime := time.Since(startTime)

		for _, entry := range lookupResults {
			fmt.Println(entry[0], "    ", entry[1])
		}

		fmt.Printf(blue+"\n>> "+reset+"%d record(s) found from your search parameters...\n", len(lookupResults))
		fmt.Printf(blue+">> "+reset+"Search completed in %s.\n", elapsedTime)
		fmt.Println()

		if len(os.Args) > 1 {
			break
		}

	}
}

func getSearchParams() []string {
	var searchParams []string

	if len(os.Args) > 1 {
		searchParams = os.Args[1:]
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if !scanner.Scan() {
				break
			}
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				break
			}
			searchParams = append(searchParams, line)
		}
	}
	fmt.Println("Searching...")
	return searchParams
}

func init() {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		reset = "\033[0m"
		red = "\033[31m"
		green = "\033[32m"
		blue = "\033[34m"
		yellow = "\033[33m"
	} else {
		reset = ""
		red = ""
		green = ""
		blue = ""
		yellow = ""
	}
}
