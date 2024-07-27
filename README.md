# fx-demo

This is a demo application that showcases how to use the Fx dependency injection and lifecycle management framework in Go.
Fx is primarily intended to be used for long-running server applications.


## Features

- **Dependency Injection:** Uses Fx to manage dependencies and improve code structure.
- **HTTP Server:** Implements a simple HTTP server using `net/http`.
- **Routing:** Defines routes using a custom `Route` interface and `ServeMux`.
- **Logging:** Utilizes `log/slog` for structured logging.
- **Graceful Shutdown:** Handles graceful shutdown of the HTTP server.

## Running the Application

1. Make sure you have Go installed on your system.
2. Navigate to the project directory.
3. Run the following command to build and run the application:

```bash
go run github.com/rezkam/fx-demo/cmd/server
```

## Code Structure

- `cmd/server/main.go`: Entry point of the application. Initializes the Fx application and registers components.
- `echo/echo.go`: Contains the `EchoHandler` that handles requests to the `/echo` route.
- `hello/hello.go`: Contains the `HelloHandler` that handles requests to the `/hello` route.
- `route/routes.go`: Defines the `Route` interface and provides a function to create a `ServeMux` with registered routes.

## Concepts Demonstrated

- **Fx Application:** Creating and running an Fx application using `fx.New()`.
- **Providers:** Defining constructors for dependencies using `fx.Provide()`.
- **Invokers:** Executing functions that depend on provided values using `fx.Invoke()`.
- **Lifecycle Hooks:** Registering functions to be executed on application startup and shutdown using `fx.Hook`.
- **Value Groups:** Grouping values of the same type using `fx.Annotate()` and `fx.ParamTags()`.
- **Result Tags:** Annotating constructors to provide values to specific groups using `fx.ResultTags()`.


## The Application Lifecycle

The application lifecycle has two main phases: **initialization** and **execution**. Both phases consist of multiple steps.

### Initialization:

1. Register all constructors passed to `fx.Provide`.
2. Register all decorators passed to `fx.Decorate`.
3. Run all functions passed to `fx.Invoke` (calling constructors and decorators as needed).

### Execution:

1. Run all hooks appended to the application by providers, decorators, or invoke functions.
2. Wait for a signal to stop running.
3. Run all shutdown hooks appended to the application.

## Lifecycle Hooks

Lifecycle hooks provide the ability to schedule work to be executed by Fx when the application starts or stops.

### Kinds of hooks:

1. **Startup hooks** (also called **OnStart** hooks): These are run in the order they are added.
2. **Shutdown hooks** (also called **OnStop** hooks): These are run in **reverse** order they were appended.

### Hooks Notes:

*  Must not block to run long-running tasks synchronously.
*  Should schedule long-running tasks to run in background goroutines.
*  Shutdown hooks should stop the background work started by the startup hooks. 

## Value Group

A value group is a collection of values that are all the same type. Any number of constructors across an Fx application can feed values to a group. Similarly, any number of consumers can consume values from a group.

**Note:** The order of values in a group is not guaranteed, so do not rely on it.


## My learning and thoughts
**Strengths:**

* **Reflection-based:** Fx leverages reflection at runtime to wire dependencies. This makes it easier to get started with and requires less boilerplate code.
* **Dynamic:** Fx allows for more dynamic behavior, enabling you to provide dependencies on the fly and modify the dependency graph at runtime.
* **Lifecycle Management:** Offers robust lifecycle management features, allowing you to define startup and shutdown logic for your components.
* **Extensible:** Provides hooks and plugins for extending its functionality.
* **Good for large, complex applications:** Its flexibility and dynamic nature make it suitable for large projects with evolving requirements.

**Weaknesses:**

* **Reflection Overhead:** Using reflection can introduce some runtime overhead compared to code generation approaches.
* **Less Compile-Time Safety:** Since wiring happens at runtime, errors in dependency configuration might only surface during execution.
* **Steeper Learning Curve:** The dynamic and flexible nature might be overwhelming for simpler projects and require more time to master.
