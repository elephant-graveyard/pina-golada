# Pina-Golada /pi:nɑ:-goʊlɑ:dɑ:/

[![License](https://img.shields.io/github/license/homeport/pina-golada.svg)](https://github.com/homeport/pina-golada/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/homeport/pina-golada)](https://goreportcard.com/report/github.com/homeport/pina-golada)
[![Build Status](https://travis-ci.org/homeport/pina-golada.svg?branch=develop)](https://travis-ci.org/homeport/pina-golada)
[![GoDoc](https://godoc.org/github.com/homeport/pina-golada?status.svg)](https://godoc.org/github.com/homeport/pina-golada)
[![Release](https://img.shields.io/github/release/homeport/pina-golada.svg)](https://github.com/homeport/pina-golada/releases/latest)

![pina-golada](.docs/logo-small.png?raw=true "Pina-Golada logo by Michael Schepanske")

## Introducing the pina-golada

Pina-Golada is a tool and framework to package assets (files and directories) into an executable binary. To be precise, it is for executables that are created by the Go compiler.

For this to work, `pina-golada` scans your Go source files for well known annotations and creates Go source code that contains the referenced assets. In your program, you are able to access the assets in memory or write them to disk by simplify defining an easy interface to work against. Pina-Golada will take care of the rest to provide the files and directories in the final build, without you having to render them in the source code yourself.

Tools that are known to use `pina-golada` to package assets into their program code are:

- [**watchful**](https://github.com/homeport/watchful) - tool to measure the disruption caused by a change to a Cloud Foundry environment
- [**gonut**](https://github.com/homeport/gonut) - portable tool to help you verify whether you can push a sample app to Cloud Foundry

_This project is work in progress._

## Installation

Installation options are either using Homebrew or a convenience download script.

- On macOS systems, a Homebrew tap is available to install `pina-golada`:

  ```sh
  brew install homeport/tap/pina-golada
  ```

- Use a convenience script to download the latest release to install it in a suitable location on your local machine (works for Linux and macOS systems):

  ```sh
  curl -fsL https://ibm.biz/Bd2645 | bash
  ```

## Usage example

In case your Go project needs assets packaged into your program, you can use `pina-golada` and the following code stub:

```go
package assets

import (
 "github.com/homeport/pina-golada/pkg/files"
)

// Provider is the provider instance to access the assets framework
var Provider ProviderInterface

// ProviderInterface is the interface that provides the assets
// @pgl(injector=Provider)
type ProviderInterface interface {
  
  // GetConfigFile returns an in-memory version of the asset
  // @pgl(asset=/assets/default-config.yaml&compressor=tar)
  GetConfigFile() (dir files.Directory , e error)
}
```

Prior to the build, run `pina-golada generate` to create a generated Go source file, which contains your assets. To clean up any generated file, use `pina-golada cleanup`. If you fancy more details, enable the verbose mode using `--verbose`. This works for `generate` as well as `cleanup`.
In addition to the verbose parameter, you can specify the annotation `parser` you want to use by providing the parser tag value.
Here are the currently supported `parser` values:

- `property`
    - Example: `@pgl(asset=/my/path&compressor=tar)`
- `csv`
    - Example: `@pgl(asset,/my/path;compressor,tar)`
- `build-tag`
    - Example: `+pgl asset,/my/path compressor,tar`   

## Contributing

We are happy to have other people contributing to the project. If you decide to do that, here's how to:

- get a Go development environment with version 1.11 or greater
- fork the project
- create a new branch
- make your changes
- open a PR.

Git commit messages should be meaningful and follow the rules nicely written down by [Chris Beams](https://chris.beams.io/posts/git-commit/):
> The seven rules of a great Git commit message
>
> 1. Separate subject from body with a blank line
> 1. Limit the subject line to 50 characters
> 1. Capitalize the subject line
> 1. Do not end the subject line with a period
> 1. Use the imperative mood in the subject line
> 1. Wrap the body at 72 characters
> 1. Use the body to explain what and why vs. how

### Running test cases and binaries generation

There are multiple make targets, but running `all` does everything you want in one call.

```sh
make all
```

### Test it with Linux on your macOS system

Best way is to use Docker to spin up a container:

```sh
docker run \
  --interactive \
  --tty \
  --rm \
  --volume $GOPATH/src/github.com/homeport/pina-golada:/go/src/github.com/homeport/pina-golada \
  --workdir /go/src/github.com/homeport/pina-golada \
  golang:1.11 /bin/bash
```

### Git pre-commit hooks

Add a pre-commit hook using this command in the repository directory:

```sh
cat <<EOS | cat > .git/hooks/pre-commit && chmod a+rx .git/hooks/pre-commit
#!/usr/bin/env bash

set -euo pipefail
make test

EOS
```

## License

Licensed under [MIT License](https://github.com/homeport/pina-golada/blob/master/LICENSE)
