package main

import (
	"go.uber.org/fx"
)

func main() {
	// we call fx.New() to setup the components

	app := fx.New()

	// run the application
	// fx is primarily intended to be used for long-running server applications
	app.Run()
}
