package main

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrUnknownCmd = errors.New("unknown command")
)

type Router map[string]Command

func (r Router) Handle(command Command) {
	r[command.Use()] = command
}

func (r Router) Call(args ...string) error {
	cmd, ok := r[args[0]]
	if !ok {
		return ErrUnknownCmd
	}
	if err := cmd.Run(context.TODO(), args[1:]); err != nil {
		return fmt.Errorf("error occured while performing command: %w", err)
	}
	return nil
}
