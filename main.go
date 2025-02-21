package main

import (
	"fmt"
	"os"
)

func main() {
	app := App{}
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
