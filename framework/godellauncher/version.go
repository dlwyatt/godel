// Copyright 2016 Palantir Technologies, Inc.
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

package godellauncher

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/palantir/godel/layout"
)

var Version = "unspecified"

func VersionTask() Task {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), layout.AppName, "version", Version)
			return nil
		},
	}
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	return Task{
		Name:        cmd.Use,
		Description: cmd.Short,
		RunImpl: func(t *Task, global GlobalConfig, stdout io.Writer) error {
			cmd.SetOutput(stdout)
			return cmd.Execute()
		},
	}
}
