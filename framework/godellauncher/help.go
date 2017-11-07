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
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"
)

type nameWithDesc struct {
	name, desc string
}

func helpFlagTask(tasks []Task, flags []nameWithDesc) Task {
	return Task{
		Name:        "--help",
		Description: "print help and exit",
		RunImpl: func(t *Task, global GlobalConfig, stdout io.Writer) error {
			fmt.Fprintln(stdout, helpContent(tasks, flags))
			return nil
		},
	}
}

func helpContent(tasks []Task, flags []nameWithDesc) string {
	output := itemsHelp("Usage", []nameWithDesc{
		{name: "godel [<global flags>] <command> [<args>]"},
	})

	var items []nameWithDesc
	for _, task := range tasks {
		items = append(items, nameWithDesc{
			name: task.Name,
			desc: task.Description,
		})
	}
	cmdsStr := itemsHelp("Commands", items)
	if cmdsStr != "" {
		output = output + "\n" + cmdsStr
	}

	flagsStr := itemsHelp("Global flags", flags)
	if flagsStr != "" {
		output = output + "\n" + flagsStr
	}
	return output
}

func itemsHelp(title string, items []nameWithDesc) string {
	if len(items) == 0 {
		return ""
	}
	output := &bytes.Buffer{}
	fmt.Fprintf(output, "%s:\n", title)

	tw := tabwriter.NewWriter(output, 0, 3, 3, ' ', 0)
	for _, item := range items {
		if item.desc == "" {
			fmt.Fprintf(tw, "\t%s\n", item.name)
		} else {
			fmt.Fprintf(tw, "\t%s\t%s\n", item.name, item.desc)
		}
	}
	_ = tw.Flush()

	return output.String()
}
