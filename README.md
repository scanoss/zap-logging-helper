# SCANOSS Platform 2.0 Zap Logging Helper Package
Welcome to the SCANOSS Platform 2.0 zap logging helper package.

This package contains helper functions to make development of Go gRPC services easier to configure for logging.

## Repository Structure
This repository is made up of the following components:
* Zap logging initialisation
* gRPC interceptor to add request/response IDs to logging

## Usage
### Zap Logger
When leveraging a zap logger in an application, it can be problematic to instantiate an instance and maintain access to this.
The zap `logger` package simplifies this by providing two global variables which give access to both the standard (faster) `Logger` and the more format friendly `Sugared Logger`.

Instantiation of these loggers is further simplified by providing opinionated versions. These include:
* Simple Dev Logger (`NewDevLogger`)
* Simple Prod Logger (`NewProdLogger`)
* Simple Dev/Prod Logger with initial Level (`NewDevLogger(level)`, `NewProdLogger(level)`)

There is also an advanced version which allowed for the importation of config from a file: `NewLoggerFromFile`

This package also provides support for dynamic level setting (`AtomicLevel`) while the application is running.
This can (optionally) be exposed to HTTP to provide external manipulation of the logging level: `SetupDynamicLogging(addr)`

### gRPC Context Server Interceptor
When working with gRPC services, it's important to provide context to all requests to aid tracing/debugging/etc.

The `github.com/grpc-ecosystem/go-grpc-middleware` package provides some amazing features to help in this area (more details [here](https://github.com/grpc-ecosystem/go-grpc-middleware)).
It provides examples and support for gathering context from a gRPC request, and interceptors to simplify this process.

The interceptors in this package provide an opinionated implementation of these examples.
Principally, these interceptors are designed to do the following:
* Check the incoming request headers for the `x-request-id` field. If none is detected, automatically generate one.
* Add this request id to all subsequent downstream calls
* Copy the `x-request-id` to the `x-response-id` and add it to the response header
* Add the `x-request-id` to the zap logging context, so that it is logged with all events generated

## Bugs/Features
To request features or alert about bugs, please do so [here](https://github.com/scanoss/zap-logging-helper/issues).

## Changelog
Details of major changes to the library can be found in [CHANGELOG.md](CHANGELOG.md).
