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
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "lintool",
	Short:             "lintool is a lint tool",
	Example:           "lintool blankline",
	SilenceErrors:     true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func main() {
	rootCmd.AddCommand(blankLineCommand(), licenseCommand(), signCommitCommand())
	if err := rootCmd.Execute(); err != nil {
		slog.Error("lintool execute failed", "reason", err)
		os.Exit(1)
	}
}
