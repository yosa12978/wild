package main

import "context"

type GenerateCommand struct {
}

func NewGenerateCmd() Command {
	return &GenerateCommand{}
}

func (g *GenerateCommand) Run(ctx context.Context, args []string) error {
	return nil
}

func (g *GenerateCommand) Long() string {
	return "this command generates a static website"
}
func (g *GenerateCommand) Short() string {
	return "generate static website"
}

func (g *GenerateCommand) Subcommands() map[string]Command {
	return nil
}

func (g *GenerateCommand) Use() string {
	return "generate"
}
