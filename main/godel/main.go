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

package main

import (
	"fmt"
	"os"

	"github.com/kardianos/osext"
	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/pkg/errors"

	"github.com/palantir/godel/framework/builtintask"
	"github.com/palantir/godel/framework/godelcli"
	"github.com/palantir/godel/framework/godellauncher"
)

func main() {
	gödelPath, err := osext.Executable()
	if err != nil {
		fmt.Printf("%+v\n", errors.Wrapf(err, "failed to determine path for current executable"))
		os.Exit(1)
	}

	if err := dirs.SetGoEnvVariables(); err != nil {
		fmt.Printf("%+v\n", errors.Wrapf(err, "failed to set Go environment variables"))
		os.Exit(1)
	}

	cmdLib, err := builtintask.AmalgomatedCmdLib(gödelPath)
	if err != nil {
		fmt.Printf("%+v\n", errors.Wrapf(err, "failed to create amalgomated CmdLib"))
		os.Exit(1)
	}
	os.Exit(amalgomated.RunApp(os.Args, nil, cmdLib, runGodelApp))
}

func runGodelApp(osArgs []string) int {
	os.Args = osArgs

	global, err := godellauncher.ParseAppArgs(os.Args)
	if err != nil {
		printErrAndExit(err, global.Debug)
	}

	var allTasks []godellauncher.Task
	allTasks = append(allTasks, godellauncher.VersionTask())
	allTasks = append(allTasks, godelcli.Tasks(global.Wrapper)...)
	allTasks = append(allTasks, builtintask.Tasks()...)
	allTasks = append(allTasks, godelcli.VerifyTask(allTasks))

	task, err := godellauncher.TaskForInput(global, allTasks)
	if err != nil {
		printErrAndExit(err, global.Debug)
	}

	if err := task.Run(global, os.Stdout); err != nil {
		printErrAndExit(err, global.Debug)
	}
	return 0
}

func printErrAndExit(err error, debug bool) {
	if errStr := err.Error(); errStr != "" {
		if debug {
			errStr = errorstringer.StackWithInterleavedMessages(err)
		}
		fmt.Fprintln(os.Stderr, errStr)
	}
	os.Exit(1)
}
