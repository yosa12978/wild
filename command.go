package main

import (
	"context"
)

type Command interface {
	Run(ctx context.Context, args []string) error

	Subcommands() map[string]Command
	Use() string
	Short() string
	Long() string
}
