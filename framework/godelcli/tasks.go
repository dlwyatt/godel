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

package godelcli

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/godellauncher"
)

func Tasks(wrapperPath string) []godellauncher.Task {
	return []godellauncher.Task{
		InstallTask(),
		UpdateTask(wrapperPath),
		CheckPathTask(),
		GitHooksTask(),
		GitHubWikiTask(),
		IDEATask(),
		PackagesTask(),
	}
}

// Creates a new godellauncher.Task that runs the provided command. The runner for the task sets stdout and sets os.Args
// to: [executable] [--wrapper <wrapper>] [task] [task args...] and invokes cmd.Execute().
func gödelCLITask(cmd *cobra.Command) godellauncher.Task {
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	return godellauncher.Task{
		Name:        cmd.Use,
		Description: cmd.Short,
		RunImpl: func(t *godellauncher.Task, global godellauncher.GlobalConfig, stdout io.Writer) error {
			cmd.SetOutput(stdout)
			args := []string{global.Executable}
			args = append(args, global.Task)
			args = append(args, global.TaskArgs...)
			os.Args = args
			return cmd.Execute()
		},
	}
}
