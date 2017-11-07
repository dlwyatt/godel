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

package builtintask

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/palantir/checks/gocd/cmd/gocd"
	"github.com/palantir/checks/gogenerate/cmd/gogenerate"
	"github.com/palantir/checks/golicense/cmd/golicense"
	"github.com/palantir/pkg/cli"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/verifyorder"
)

func builtinTasks() []godellauncher.Task {
	return []godellauncher.Task{
		createBuiltinVerifyTask("generate", gogenerate.App(), "generate.yml", verifyorder.Generate),
		createBuiltinVerifyTask("imports", gocd.App(), "imports.yml", verifyorder.Imports),
		createBuiltinVerifyTask("license", golicense.App(), "license.yml", verifyorder.License),
		createDisgtoTask("run"),
		createDisgtoTask("project-version"),
		createDisgtoTask("build"),
		createDisgtoTask("dist"),
		createDisgtoTask("clean"),
		createDisgtoTask("artifacts"),
		createDisgtoTask("products"),
		createDisgtoTask("docker"),
		createDisgtoTask("publish"),
	}
}

func createDisgtoTask(name string) godellauncher.Task {
	task := createBuiltinTaskHelper(name, distgo.App(), []string{name}, "dist.yml", false, 0)

	baseRunImpl := task.RunImpl
	// distgo requires working directory to be the base (project) directory, so decorate action to set the working
	// directory before invocation.
	task.RunImpl = func(t *godellauncher.Task, global godellauncher.GlobalConfig, stdout io.Writer) error {
		if global.Wrapper != "" {
			if !filepath.IsAbs(global.Wrapper) {
				absWrapperPath, err := filepath.Abs(global.Wrapper)
				if err != nil {
					return errors.Wrapf(err, "failed to convert wrapper path to absolute path")
				}
				global.Wrapper = absWrapperPath
			}
			global.WorkingDir = path.Dir(global.Wrapper)
			if err := os.Chdir(global.WorkingDir); err != nil {
				return errors.Wrapf(err, "failed to change working directory")
			}
		}
		return baseRunImpl(t, global, stdout)
	}
	return task
}

func createBuiltinVerifyTask(name string, app *cli.App, cfgFileName string, verifyOrder int) godellauncher.Task {
	return createBuiltinTaskHelper(name, app, nil, cfgFileName, true, verifyOrder)
}

func createBuiltinTaskHelper(name string, app *cli.App, cmdPath []string, cfgFileName string, verify bool, verifyOrder int) godellauncher.Task {
	currCmd := app.Command
	for _, wantSubCmdName := range cmdPath {
		for _, currSubCmd := range currCmd.Subcommands {
			if currSubCmd.Name == wantSubCmdName {
				currCmd = currSubCmd
				break
			}
		}
	}
	var verifyTask *godellauncher.VerifyOptions
	if verify {
		verifyTask = &godellauncher.VerifyOptions{
			Ordering: verifyOrder,
		}
	}

	return godellauncher.Task{
		Name:        name,
		Description: currCmd.Usage,
		ConfigFile:  cfgFileName,
		Verify:      verifyTask,
		RunImpl: func(t *godellauncher.Task, global godellauncher.GlobalConfig, stdout io.Writer) error {
			args, err := CfgCLIArgs(global, cmdPath, t.ConfigFile)
			if err != nil {
				return err
			}
			app.Stdout = stdout
			app.Stderr = os.Stderr
			os.Args = args
			if exitCode := app.Run(args); exitCode != 0 {
				return fmt.Errorf("")
			}
			return nil
		},
	}
}
