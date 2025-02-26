package main

import (
	"os"
)

type App struct {
}

func (a *App) Run() error {
	draftCmd := &DraftCommand{}
	versionCmd := &VersionCommand{}
	generateCmd := &GenerateCommand{}

	router := Router{}
	router.Handle(draftCmd)
	router.Handle(versionCmd)
	router.Handle(generateCmd)
	return router.Call(os.Args[1:]...)
}
