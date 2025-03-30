// Copyright 2025 openHoles Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"path/filepath"
	"strings"
)

// LicenseType represents the type of open source license
type LicenseType string

const (
	Apache2    LicenseType = "Apache-2.0"
	MIT        LicenseType = "MIT"
	GPL2       LicenseType = "GPL-2.0"
	GPL3       LicenseType = "GPL-3.0"
	LGPL       LicenseType = "LGPL"
	Mozilla    LicenseType = "MPL"
	BSD        LicenseType = "BSD"
	Unlicensed LicenseType = "UNLICENSED"
)

// licensePatterns defines characteristic text patterns for different licenses
var licensePatterns = map[LicenseType][]string{
	Apache2: {"Apache License", "Version 2.0", "http://www.apache.org/licenses/LICENSE-2.0"},
	MIT:     {"MIT License", "Permission is hereby granted, free of charge"},
	GPL2:    {"GNU GENERAL PUBLIC LICENSE", "Version 2"},
	GPL3:    {"GNU GENERAL PUBLIC LICENSE", "Version 3"},
	LGPL:    {"GNU LESSER GENERAL PUBLIC LICENSE"},
	Mozilla: {"Mozilla Public License"},
	BSD:     {"BSD License", "Redistribution and use in source and binary forms"},
}

// detectLicenseType detects the type of license in the project
func detectLicenseType(directory string) (LicenseType, error) {
	// Check common LICENSE file names
	licenseFiles := []string{
		"LICENSE",
		"LICENSE.txt",
		"LICENSE.md",
		"COPYING",
		"COPYING.txt",
	}

	var licenseContent string
	for _, name := range licenseFiles {
		path := filepath.Join(directory, name)
		content, err := os.ReadFile(path)
		if err == nil {
			licenseContent = string(content)
			break
		}
	}

	// If no LICENSE file is found, default to Apache 2.0
	if licenseContent == "" {
		return Unlicensed, nil
	}

	// Detect license type
	for licType, patterns := range licensePatterns {
		matched := true
		for _, pattern := range patterns {
			if !strings.Contains(licenseContent, pattern) {
				matched = false
				break
			}
		}
		if matched {
			return licType, nil
		}
	}

	return Apache2, nil // If unable to determine, default to Apache 2.0
}

// validateLicenseHeader checks if the file contains the correct license header
func validateLicenseHeader(content string, licenseType LicenseType) bool {
	switch licenseType {
	case Apache2:
		return strings.Contains(content, "Licensed under the Apache License, Version 2.0")
	case MIT:
		return strings.Contains(content, "MIT License")
	case GPL2:
		return strings.Contains(content, "GNU GENERAL PUBLIC LICENSE") &&
			strings.Contains(content, "Version 2")
	case GPL3:
		return strings.Contains(content, "GNU GENERAL PUBLIC LICENSE") &&
			strings.Contains(content, "Version 3")
	default:
		return false
	}
}
