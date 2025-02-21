package main

import (
	"context"
	"fmt"
)

type VersionCommand struct {
}

func (v *VersionCommand) Run(ctx context.Context, args []string) error {
	fmt.Println("0.0.1")
	return nil
}

func (v *VersionCommand) Subcommands() map[string]Command {
	return nil
}

func (v *VersionCommand) Use() string {
	return "version"
}

func (v *VersionCommand) Short() string {
	return "shows version"
}

func (v *VersionCommand) Long() string {
	return "shows version long"
}
