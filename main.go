package main

import (
	"context"
	"errors"
	"github.com/rezkam/fx-demo/hello"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/rezkam/fx-demo/echo"
	"github.com/rezkam/fx-demo/route"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	// we call fx.New() to setup the components of the application.
	app := fx.New(
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),
		// provide adds a http server to the application. The server hooks into the application lifecycle.
		// so it will start serving requests when the application starts and stop when the application stops.
		fx.Provide(
			// order of the constructors given to fx.Provide does *not* matter.
			NewHTTPServer,
			fx.Annotate(
				route.NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			AsRoute(echo.NewHandler),
			// Fx does not allow two constructors to provide the same type without annotating them.
			// Here we need to annotate the NewHandler and NewHandler constructors to distinguish them.
			// using fx.ResultTag
			NewJSONLogger,
			AsRoute(hello.NewHandler),
		),
		// fx.Invoke used to request that the HTTP Server always instantiated
		// even if none of the other components in the application reference it directly.
		fx.Invoke(func(s *http.Server) {}),
	)

	// run the application
	// fx is primarily intended to be used for long-running server applications.
	app.Run()
}

func NewJSONLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, nil))
}

// NewHTTPServer builds an HTTP server that will begin serving requests
// when the Fx application starts.
func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, logger *slog.Logger) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}
	// fx hooks are functions that are executed at different points in the application lifecycle.
	// here hooks are used to start and stop the server
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			logger.Info("Starting HTTP server", "addr", srv.Addr)
			//hooks must not block to run a long-running task synchronously, so we run the server in a goroutine.
			go func() {
				if err := srv.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
					logger.Error("HTTP server stopped", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// we need to stop the background work started by the startup hooks.
			// so we call Shutdown on the server to stop it gracefully
			// without interrupting any active connections.
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

/* 	The Application Lifecycle
Has two main phases: initialization and execution, and both consist of multiple steps.
Initialization:
1. Register all constructors passed to fx.Provide.
2. Register all decorators passed to fx.Decorate.
3. Run all functions passed to fx.Invoke (calling constructors and decorators as needed).
Execution:
1. Run all hooks appended to the application by providers, decorators, or invoke functions.
2. Wait for a signal to stop running
3. Run all shutdown hooks appended to the application.

Lifecycle Hooks provide the ability to schedule work to be executed by Fx when the application starts or stops.
Kinds of hooks:
1. Startup hooks also called OnStart hooks these are run in the order they are added.
2. Shutdown hooks also called OnStop hooks these are run in reverse order they were appended.

Hooks must not block to run long-running tasks synchronously.
Hooks should schedule long-running tasks to run in the background goroutines.
Shutdown hooks should stop the background work started by the startup hooks.
*/

// Value Group
// Is a collection of values that are all the same type.
// Any number of constructors across Fx application can feed values to a group.
// Similarly, any number of consumers can consume values from a group.
// Here we defined the group "routes" to hold all the Route instances.
// Handlers are annotated with the group tag to indicate that they should be added to the group.
// The NewServeMux constructor is annotated with the ParamTags("group:routes") to indicate that it should receive all the Route instances.
// The order of values in a group is not guaranteed so do not rely on the order of values in a group.

// AsRoute annotates the given constructor to state that
// it provides a route to the "routes" group.
func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(route.Route)),
		fx.ResultTags(`group:"routes"`),
	)
}
