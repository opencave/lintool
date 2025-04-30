// Copyright 2025 opencave Authors
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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	excludeSourceCodeExtensions = []string{"^\\..*", ".mod", ".sum", ".md"} // ignore dotfile, project configuration file
	excludeFiles                = []string{"^\\..*", "Dockerfile.*", "Makefile.*"}
)

func licenseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "license",
		Short:        "check if the source code files have a license header",
		Aliases:      []string{"lic"},
		SilenceUsage: true,
		Example:      "lintool lic -d . -e .idea,testdata",
		RunE: func(cmd *cobra.Command, args []string) error {
			directory, err := cmd.Flags().GetString("directory")
			if err != nil {
				return err
			}
			exclude, err := cmd.Flags().GetStringSlice("exclude")
			if err != nil {
				return err
			}
			license, err := cmd.Flags().GetString("license")
			if err != nil {
				return err
			}
			excludeExtentions, err := cmd.Flags().GetStringSlice("exclude-extensions")
			if err != nil {
				return err
			}
			return checkLicenseHeader(directory, license, exclude, excludeExtentions)
		},
	}
	cmd.Flags().StringP("directory", "d", ".", "directory to check")
	cmd.Flags().StringSliceP("exclude", "e", excludeFiles, "directories or files to exclude (comma-separated)")
	cmd.Flags().StringP("license", "l", "", "license file to use, if no license file is found and this flag is not set, process will be skipped, support [Apache-2.0, MIT, GPL-2.0, GPL-3.0, LGPL, MPL, BSD]")
	cmd.Flags().StringSliceP("exclude-extensions", "x", excludeSourceCodeExtensions, "exclude file extension to check, if file no extension will be skipped (comma-separated)")
	return cmd
}

func checkLicenseHeader(directory, license string, exclude, excludeExtensions []string) error {
	// Detect the project's license type
	var err error
	var licenseType = Unlicensed
	if license != "" {
		licenseType = LicenseType(license)
	} else {
		licenseType, err = detectLicenseType(directory)
		if err != nil {
			return err
		}
		if licenseType == Unlicensed {
			fmt.Println("no license file found, skipping license header check")
			return errors.New("no license file found")
		}
	}

	excludeExtensionMap := make(map[string]bool)
	for _, extension := range excludeExtensions {
		excludeExtensionMap[extension] = true
	}

	excludeMap := make(map[string]bool)
	for _, e := range exclude {
		excludeMap[e] = true
	}

	var errorFiles []string
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}

		// Check excluded paths
		current := relPath
		for current != "." && current != "" {
			if excludeMap[current] {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			for _, rule := range exclude {
				compile := regexp.MustCompile(rule)
				if compile.MatchString(current) {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
			current = filepath.Dir(current)
		}

		if info.IsDir() {
			return nil
		}

		// Only check exclude extension source code files
		if ignoreExcludeExtension(path, excludeExtensionMap) {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if !validateLicenseHeader(string(b), licenseType) {
			errorFiles = append(errorFiles, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if len(errorFiles) > 0 {
		fmt.Printf("Detected license type: %s\n", licenseType)
		fmt.Printf("The following files are missing the correct license header:\n%s\n", strings.Join(errorFiles, "\n"))
		return errors.New("license header checker issue found")
	}

	return nil
}

// ignoreExcludeExtension determines if a file is a exclude extension source code file
func ignoreExcludeExtension(path string, excludeExtensions map[string]bool) bool {
	if len(excludeExtensions) == 0 {
		return false
	}
	ext := strings.ToLower(filepath.Ext(path))
	if _, exist := excludeExtensions[ext]; exist && ext != "" {
		return true
	}
	basename := filepath.Base(path)
	for extension := range excludeExtensions {
		compile := regexp.MustCompile(extension)
		if compile.MatchString(basename) {
			return true
		}
	}
	return false
}
