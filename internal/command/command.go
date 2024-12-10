// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/googleapis/generator/internal/container"
)

type Command struct {
	Name  string
	Short string
	flags *flag.FlagSet
	Run   func(ctx context.Context) error
}

func (c *Command) Parse(args []string) error {
	return c.flags.Parse(args)
}

func Lookup(name string) (*Command, error) {
	var cmd *Command
	for _, sub := range Commands {
		if sub.Name == name {
			cmd = sub
		}
	}
	if cmd == nil {
		return nil, fmt.Errorf("invalid command: %q", name)
	}
	return cmd, nil
}

var CmdConfigure = &Command{
	Name:  "configure",
	Short: "Configure a new API in a given language",
	Run: func(ctx context.Context) error {
		if flagAPIRoot == "" {
			return fmt.Errorf("-api-root is not provided")
		}
		if flagAPIPath == "" {
			return fmt.Errorf("-api-path is not provided")
		}
		if flagPush && flagGitHubToken == "" {
			return fmt.Errorf("-github-token must be provided if -push is set to true")
		}
		return container.Configure(ctx, flagAPIRoot, flagAPIPath, "")
	},
}

var CmdGenerate = &Command{
	Name:  "generate",
	Short: "Generate a new client library",
	Run: func(ctx context.Context) error {
		if flagAPIRoot == "" {
			return fmt.Errorf("-api-root is not provided")
		}
		if flagAPIPath == "" {
			return fmt.Errorf("-api-path is not provided")
		}
		return container.Generate(ctx, flagAPIRoot, flagAPIPath, flagOutput, "")
	},
}

var CmdUpdateRepo = &Command{
	Name:  "update-repo",
	Short: "Configure a new API in a given language",
	Run: func(ctx context.Context) error {
		if err := container.Generate(ctx, flagAPIRoot, flagAPIPath, flagOutput, ""); err != nil {
			return err
		}
		if err := container.Clean(ctx, flagOutput, flagAPIPath); err != nil {
			return err
		}
		if err := container.Build(ctx, flagOutput, flagAPIPath); err != nil {
			return err
		}
		if err := commit(); err != nil {
			return err
		}
		return push()
	},
}

func commit() error {
	return nil
}

func push() error {
	return nil
}

var Commands = []*Command{
	CmdConfigure,
	CmdGenerate,
	CmdUpdateRepo,
}

func init() {
	for _, c := range Commands {
		c.flags = flag.NewFlagSet(c.Name, flag.ContinueOnError)
		c.flags.Usage = constructUsage(c.flags, c.Name)
	}

	fs := CmdConfigure.flags
	for _, fn := range []func(fs *flag.FlagSet){
		addFlagAPIPath,
		addFlagAPIRoot,
		addFlagLanguage,
		addFlagPush,
	} {
		fn(fs)
	}

	fs = CmdGenerate.flags
	for _, fn := range []func(fs *flag.FlagSet){
		addFlagAPIPath,
		addFlagAPIRoot,
		addFlagLanguage,
		addFlagOutput,
	} {
		fn(fs)
	}

	fs = CmdUpdateRepo.flags
	for _, fn := range []func(fs *flag.FlagSet){
		addFlagAPIPath,
		addFlagAPIRoot,
		addFlagBranch,
		addFlagGitHubToken,
		addFlagLanguage,
		addFlagOutput,
		addFlagPush,
	} {
		fn(fs)
	}
}

func constructUsage(fs *flag.FlagSet, name string) func() {
	output := fmt.Sprintf("Usage:\n\n  generator %s [arguments]\n", name)
	output += "\nFlags:\n\n"
	return func() {
		fmt.Fprint(fs.Output(), output)
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\n\n")
	}
}
