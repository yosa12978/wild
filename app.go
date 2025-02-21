package main

import (
	"os"
)

type App struct {
}

func (a *App) Run() error {
	draftCmd := &DraftCommand{}
	versionCmd := &VersionCommand{}

	router := Router{}
	router.Handle(draftCmd)
	router.Handle(versionCmd)
	return router.Call(os.Args[1:]...)
}
