package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/fx"
)

func main() {
	// we call fx.New() to setup the components of the application.
	app := fx.New(
		// provide adds a http server to the application. The server hooks into the application lifecycle.
		// so it will start serving requests when the application starts and stop when the application stops.
		fx.Provide(NewHTTPServer),
		// fx.Invoke used to request that the HTTP Server always instantiated
		// even if none of the other components in the application reference it directly.
		fx.Invoke(func(s *http.Server) {}),
	)

	// run the application
	// fx is primarily intended to be used for long-running server applications.
	app.Run()
}

// NewHTTPServer builds an HTTP server that will begin serving requests
// when the Fx application starts.
func NewHTTPServer(lc fx.Lifecycle) *http.Server {
	srv := &http.Server{Addr: ":8080"}
	// fx hooks are functions that are executed at different points in the application lifecycle.
	// here hooks are used to start and stop the server
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Printf("Server listening on %s", srv.Addr)
			//hooks must not block to run a long-running task synchronously so we run the server in a goroutine.
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

/* 	Aplication Lifecycle
Has two main phases: initialization and execution and both are comprised of multiple steps.
Initialization:
1. Register all constructors passed to fx.Provide.
2. Register all decorators passed to fx.Decorate.
3. Run all functions passed to fx.Invoke (calliing constructors and decorators as needed).
Execution:
1. Run all hooks appended to the application by providers, decorators, or invoke functions.
2. Wait for a signal to stop running
3. Run all shutdown hooks appended to the application.

Lifecycle Hooks
Lifecycle hooks provide the ability to schedule work to be executed by Fx when the application starts or stops.
Kinds of hooks:
1. Startup hooks also called OnStart hooks these are run in the order they are added.
2. Shutdown hooks also called OnStop hooks these are run in reverse order they were appended.

Hooks must not block to run long-running tasks synchronously.
hooks should schedule long-running tasks to run in the background goroutines.
shutdown hooks should stop the background work started by the startup hooks.
*/
