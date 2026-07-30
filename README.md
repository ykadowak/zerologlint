# zerologlint
![build](https://github.com/ykadowak/zerologlint/actions/workflows/testing.yaml/badge.svg)

`zerologlint` is a linter for [zerolog](https://github.com/rs/zerolog) that can be run with `go vet` or through [golangci-lint](https://golangci-lint.run/) since `v1.53.0`.
It detects the wrong usage of `zerolog` that a user forgets to dispatch `zerolog.Event` with `Send` or `Msg` like functions, in which case nothing will be logged. For more detailed explanations of the cases it detects, see [Examples](#Examples).

## Install

```bash
go install github.com/ykadowak/zerologlint/cmd/zerologlint@latest
```

## Usage
```bash
go vet -vettool=`which zerologlint` ./...
```

or you can also use it with [golangci-lint](https://golangci-lint.run/) since `v1.53.0`.

## Examples
```go
package main

import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    // 1. Basic case
    log.Info() // "must be dispatched by Msg or Send method"

    // 2. Nested case
    log.Info(). // "must be dispatched by Msg or Send method"
        Str("foo", "bar").
        Dict("dict", zerolog.Dict().
            Str("bar", "baz").
            Int("n", 1),
        )

    // 3. Reassignment case
    logger := log.Info() // "must be dispatched by Msg or Send method"
    if err != nil {
        logger = log.Error() // "must be dispatched by Msg or Send method"
    }
    logger.Str("foo", "bar")

    // 4. Deferred case
    defer log.Info() // "must be dispatched by Msg or Send method"

    // 5. zerolog.Logger case
    logger2 := zerolog.New(os.Stdout)
    logger2.Info().Send()

    // 6. Dispatch in other function case
    event := log.Info()
    dispatcher(event)
}

func dispatcher(e *zerolog.Event) {
    e.Send()
}
```

## Custom Module Paths (zerolog forks/mirrors)

By default, `zerologlint` only inspects code that uses `github.com/rs/zerolog`.
If your project uses a fork or mirror of zerolog under a different module path, you can configure
additional module path prefixes.

### Using `go vet` / CLI flag

Pass the `-zerologlint.prefix` flag with a comma-separated list of additional prefixes:

```bash
go vet -vettool=`which zerologlint` -zerologlint.prefix=myorg/myzerolog ./...
# multiple prefixes:
go vet -vettool=`which zerologlint` -zerologlint.prefix=myorg/myzerolog,other/zerolog ./...
```

### Using golangci-lint (plugin mode)

When building a golangci-lint plugin, pass flags via `ldflags`:

```bash
go build -buildmode=plugin \
  -ldflags "-X 'main.flags=-prefix myorg/myzerolog'" \
  github.com/ykadowak/zerologlint/plugin/zerologlint
```

### Programmatic API (golangci-lint custom linter integrations)

Use `NewAnalyzerForSettings` to create an analyzer with settings provided programmatically:

```go
import "github.com/ykadowak/zerologlint"

analyzer := zerologlint.NewAnalyzerForSettings(zerologlint.Settings{
    AdditionalPrefixes: []string{"myorg/myzerolog"},
})
```
