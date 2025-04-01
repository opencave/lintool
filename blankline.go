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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// blankLineCommand creates a cobra.Command for checking if files end with a blank line.
//
// Flags:
//
//	-d, --directory string   directory to check (default ".")
//	-e, --exclude strings    directories or files to exclude (comma-separated)
func blankLineCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "blankline",
		Short:        "check if file ends with a blank line",
		Aliases:      []string{"bl"},
		SilenceUsage: true,
		Example:      "lintool bl -d . -e .idea,testdata",
		RunE: func(cmd *cobra.Command, args []string) error {
			directory, err := cmd.Flags().GetString("directory")
			if err != nil {
				return err
			}
			exclude, err := cmd.Flags().GetStringSlice("exclude")
			if err != nil {
				return err
			}
			return checkBlankLine(directory, exclude)
		},
	}
	cmd.Flags().StringP("directory", "d", ".", "directory to check")
	cmd.Flags().StringSliceP("exclude", "e", []string{"^\\..*", "pb.go$", "gen.go$"}, "directories or files to exclude (comma-separated)")
	return cmd
}

// checkBlankLine walks through the directory and checks if each file ends with a blank line.
// It uses the exclude list to filter out directories or files.
// If any file does not end with a blank line, it prints the file path and returns an error.
func checkBlankLine(directory string, exclude []string) error {
	excludeMap := make(map[string]bool)
	for _, e := range exclude {
		excludeMap[e] = true
	}

	var errorFiles []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the path or any of its parent directories are in the exclusion list
		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}

		// Check if the current path or any of its parent directories should be excluded
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

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if len(b) == 0 || b[len(b)-1] != '\n' {
			errorFiles = append(errorFiles, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(errorFiles) > 0 {
		fmt.Printf("the following files are not end with a blank line:\n%s\n", strings.Join(errorFiles, "\n"))
		return errors.New("blank line issue found")
	}
	return nil
}
