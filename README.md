# Actual Sync

[![Go](https://github.com/nathanjisaac/actual-sync/actions/workflows/go.yml/badge.svg)](https://github.com/nathanjisaac/actual-sync/actions/workflows/go.yml)

## ⚠️⚠️ Work In Progress ⚠️⚠️

The goal for this project is to migrate the [actualbudget/actual-server](https://github.com/actualbudget/actual-server)
source code over to use [Go](https://go.dev/).

## Roadmap

Current roadmap is documented on the [roadmap discussion](https://github.com/nathanjisaac/actual-sync/discussions/1).

## Architecture

The architecture is still somewhat of a WIP. But the high level gist looks like this.

```text
main.go         - Root Entrypoint
cmd/            - CLI Entrypoint
internal/
    server.go   - Echo server implementation
    core        - Domain logic handler functions to be used in the routes
    routes      - Echo route handlers
    storage     - Implementations for the storage gateway [sqlite, memory] - (PostgreSQL) to follow if needed
```

## CLI Usage

### actual-sync serve

This command will start the actual-sync server

#### Synopsis

This command will start the actual-sync server with the
specified configurations along with this command.

```shell
actual-sync serve [flags]
```

#### Options

```text
  -d, --data-path string   Sets data directory path (default "data")
  -l, --headless           Runs actual-sync without the web app
  -h, --help               help for serve
  -p, --port int           Runs actual-sync at specified port (default 5006)
      --production         Runs actual-sync in production mode
```

#### Global options

```text
      --config string   config file (default is $HOME/.actual-sync.yaml)
```

## Development

### Dependencies

#### Build dependencies

- [Node.js](https://nodejs.dev/)
- [go](https://go.dev/)

#### Development dependencies

- [golangci-lint](https://golangci-lint.run/)

### Steps to run

1. Install node_modules (actual-web)

```shell
$ npm install
```

2. Run go program

```shell
$ go run main.go serve
```

**NOTE: Run `golangci-lint run` after changes(if any) to ensure code quality.**
