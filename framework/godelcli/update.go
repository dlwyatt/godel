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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/godelcli/installupdate"
	"github.com/palantir/godel/framework/godellauncher"
)

func UpdateTask(wrapperPath string) godellauncher.Task {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Download and install the version of gödel specified in the godel.properties file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wrapperPath == "" {
				return errors.Errorf("wrapper path not specified")
			}
			return installupdate.Update(wrapperPath, cmd.OutOrStdout())
		},
	}
	return gödelCLITask(cmd)
}