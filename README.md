[![Go](https://github.com/nathanjisaac/actual-sync/actions/workflows/go.yml/badge.svg)](https://github.com/nathanjisaac/actual-sync/actions/workflows/go.yml)

## ⚠️⚠️ Work In Progress ⚠️⚠️

The goal for this project is to migrate the [actualbudget/actual-server](https://github.com/actualbudget/actual-server) 
source code over to use [Go](https://go.dev/).

## Roadmap

Current roadmap is documented on the [roadmap discussion](https://github.com/nathanjisaac/actual-sync/discussions/1).

## Architecture

The architecture is still somewhat of a WIP. But the high level gist looks like this. 

```
main.go - Root Entrypoint
cmd/ - CLI Entrypoint
internal/
    core - Domain logic handler functions to be used in the routes
    routes - Echo route handlers
    storage - Impelemtnations for the storage gateway [sqlite, files] - (PostgreSQL) to follow if needed
```

## CLI Usage

> Planning to make this a standardized binary once the CLI api is stable.

### Run Server

```shell
$ go run main.go serve
```
