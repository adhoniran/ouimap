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
	"compress/gzip"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ouiInfo struct {
	ouiPrefix    uint64
	prefixLength int
	vendorName   string
}

var ouiRecords []ouiInfo

var vendorIndex = make(map[string][]ouiInfo)

func loadOUIDatabase(db string) error {
	f, err := os.Open(db)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer func() { _ = f.Close() }()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("error reading database: %w", err)
	}
	defer func() { _ = gz.Close() }()

	scanner := bufio.NewScanner(gz)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.FieldsFunc(line, func(r rune) bool { return r == '\t' })
		if len(parts) < 3 {
			continue
		}

		ouiPrefix := strings.TrimSpace(parts[0])
		vendorName := strings.TrimSpace(parts[2])

		pfx, err := parseOUIEntry(ouiPrefix)
		if err != nil {
			continue
		}
		pfx.vendorName = vendorName

		ouiRecords = append(ouiRecords, pfx)

		vendorUpper := strings.ToUpper(vendorName)
		vendorIndex[vendorUpper] = append(vendorIndex[vendorUpper], pfx)
	}
	return scanner.Err()
}

func parseOUIEntry(entry string) (ouiInfo, error) {
	var result ouiInfo

	delimiterIndex := strings.IndexRune(entry, '/')
	prefixBits := ""
	if delimiterIndex != -1 {
		prefixBits = entry[delimiterIndex+1:]
		entry = entry[:delimiterIndex]
	}

	clean := strings.ReplaceAll(entry, ":", "")
	clean = strings.ReplaceAll(clean, "-", "")
	clean = strings.ReplaceAll(clean, ".", "")
	clean = strings.ToUpper(clean)

	if len(clean) != 6 && len(clean) != 8 && len(clean) != 10 && len(clean) != 12 {
		return result, fmt.Errorf("invalid hex lenghth: %s", clean)
	}

	ouiNumeric, err := strconv.ParseUint(clean, 16, 64)
	if err != nil {
		return result, fmt.Errorf("failed to parse: %w", err)
	}

	bitCount := len(clean) * 4
	remainingBits := 48 - bitCount
	ouiNumeric <<= remainingBits

	if prefixBits == "" {
		result.prefixLength = bitCount
	} else {
		m, err := strconv.Atoi(prefixBits)
		if err != nil || m < 0 || m > 48 {
			return result, fmt.Errorf("invalid mask: %s", prefixBits)
		}
		result.prefixLength = m
	}

	bitShiftCount := 48 - result.prefixLength
	shiftedOUI := (ouiNumeric >> bitShiftCount) << bitShiftCount

	result.ouiPrefix = shiftedOUI
	return result, nil
}

func promptSearchParams() []string {
	var searchParams []string
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

	fmt.Println("Searching...")
	return searchParams
}

func lookupOUIData(searchParams []string) [][]string {
	var results [][]string

	for _, entry := range searchParams {
		if matches := searchOUIEntries(entry); len(matches) > 0 {
			results = append(results, matches...)
			continue
		}

		if partialHexEntry(entry) {
			continue
		}

		if vendorMatches := searchVendorEntries(entry); len(vendorMatches) > 0 {
			results = append(results, vendorMatches...)
		}
	}

	fmt.Println(yellow)
	return results
}

func searchOUIEntries(entry string) [][]string {
	var matches [][]string
	parsed, err := parseOUIEntry(entry)
	if err != nil {
		return matches
	}

	foundMap := make(map[string]bool)
	for _, record := range ouiRecords {
		if !compareOUIPrefixes(parsed.ouiPrefix, parsed.prefixLength, record.ouiPrefix, record.prefixLength) {
			continue
		}

		uniqueKey := fmt.Sprintf("%012X/%d", record.ouiPrefix, record.prefixLength)
		if foundMap[uniqueKey] {
			continue
		}

		matches = append(matches, []string{
			formatMacPrefix(record.ouiPrefix, record.prefixLength),
			record.vendorName,
		})
		foundMap[uniqueKey] = true
	}

	return matches
}

func partialHexEntry(entry string) bool {
	clean := strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(entry, ":", ""),
			"-", ""),
		".", "")
	return len(clean) >= 2 && len(clean) < 6 && validateHex(clean)
}

func searchVendorEntries(entry string) [][]string {
	var matches [][]string
	vendorMap := make(map[string]bool)
	upperEntry := strings.ToUpper(entry)

	for vend, plist := range vendorIndex {
		if !strings.Contains(vend, upperEntry) {
			continue
		}

		for _, p := range plist {
			uniqueKey := fmt.Sprintf("%012X/%d", p.ouiPrefix, p.prefixLength)
			if vendorMap[uniqueKey] {
				continue
			}

			matches = append(matches, []string{
				formatMacPrefix(p.ouiPrefix, p.prefixLength),
				p.vendorName,
			})
			vendorMap[uniqueKey] = true
		}
	}

	return matches
}
func validateHex(inputString string) bool {
	for i := 0; i < len(inputString); i++ {
		c := inputString[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

func compareOUIPrefixes(uPrefix uint64, uLen int, bPrefix uint64, bLen int) bool {

	lMin := uLen
	if bLen < lMin {
		lMin = bLen
	}

	shift := 48 - lMin
	truncatedU := (uPrefix >> shift) << shift
	truncatedB := (bPrefix >> shift) << shift

	return truncatedU == truncatedB
}

func formatMacPrefix(prefix uint64, length int) string {
	fmtPrefix := fmt.Sprintf("%012X", prefix)
	mac := fmtPrefix[0:2] + ":" + fmtPrefix[2:4] + ":" + fmtPrefix[4:6] + ":" + fmtPrefix[6:8] + ":" + fmtPrefix[8:10] + ":" + fmtPrefix[10:12]
	if length < 48 {
		mac += fmt.Sprintf("/%d", length)
	}
	return mac
}

func deduplicateInput(searchParams []string) []string {
	dedupMap := make(map[string]bool)
	var unique []string

	for _, entry := range searchParams {
		upperTerm := strings.ToUpper(entry)

		if !dedupMap[upperTerm] {
			dedupMap[upperTerm] = true
			unique = append(unique, entry)
		}
	}

	return unique
}

func deduplicateResults(records [][]string) [][]string {
	seen := make(map[string]bool)
	var unique [][]string

	for _, item := range records {
		key := item[0] + "||" + item[1]
		if !seen[key] {
			seen[key] = true
			unique = append(unique, item)
		}
	}
	return unique
}
