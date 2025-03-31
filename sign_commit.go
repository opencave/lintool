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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func signCommitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sign_commit",
		Aliases: []string{"sc"},
		Short:   "check if commit is signed",
		Example: "lintool sign_commit",
		RunE: func(cmd *cobra.Command, args []string) error {
			directory, err := cmd.Flags().GetString("directory")
			if err != nil {
				return err
			}
			debug, err := cmd.Flags().GetBool("debug")
			if err != nil {
				return err
			}
			repository, err := git.PlainOpen(directory)
			if err != nil {
				return err
			}
			reference, err := repository.Head()
			if err != nil {
				return err
			}
			commit, err := repository.CommitObject(reference.Hash())
			if err != nil {
				return err
			}
			if commit.PGPSignature == "" {
				return errors.New("no PGP signature found in commit")
			}
			if debug {
				commitMessage, err := json.Marshal(commit)
				if err != nil {
					return err
				}
				fmt.Println(string(commitMessage))
			}
			return nil
		},
	}
	cmd.Flags().StringP("directory", "d", ".", "repository directory")
	cmd.Flags().BoolP("debug", "D", false, "show commit information")
	return cmd
}
