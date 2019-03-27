# Pina-Golada /pi:nɑ:-goʊlɑ:dɑ:/

<img src="/docs/logo.png" align="right" width="120" height="178" title="Pina-Golada logo by Michael Schepanske"/>

[![License](https://img.shields.io/github/license/homeport/pina-golada.svg)](https://github.com/homeport/pina-golada/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/homeport/pina-golada)](https://goreportcard.com/report/github.com/homeport/pina-golada)
[![Build Status](https://travis-ci.org/homeport/pina-golada.svg?branch=develop)](https://travis-ci.org/homeport/pina-golada)
[![GoDoc](https://godoc.org/github.com/homeport/pina-golada?status.svg)](https://godoc.org/github.com/homeport/pina-golada)

## Introducing the pina-golada

Pina-Golada is a tool to automatically generate asset providers that implement custom interfaces.
The tool's main propose is to hide ugly and blown up asset provider behind interfaces.  

_This project is work in progress._

## Contributing

We are happy to have other people contributing to the project. If you decide to do that, here's how to:

- get a Go development environment with version 1.11 or greater
- fork the project
- create a new branch
- make your changes
- open a PR.

Git commit messages should be meaningful and follow the rules nicely written down by [Chris Beams](https://chris.beams.io/posts/git-commit/):
> The seven rules of a great Git commit message
> 1. Separate subject from body with a blank line
> 1. Limit the subject line to 50 characters
> 1. Capitalize the subject line
> 1. Do not end the subject line with a period
> 1. Use the imperative mood in the subject line
> 1. Wrap the body at 72 characters
> 1. Use the body to explain what and why vs. how

### Installation

To install pina-golada on macOS systems, you can use the homeport Homebrew tap:

```sh
brew install homeport/tap/pina-golada
```

### Running test cases and binaries generation

There are multiple make targets, but running `all` does everything you want in one call.

```sh
make all
```

### Test it with Linux on your macOS system

Have the project that is in need of asset management depend on Pina-Golada.
Define your setup as following: 

```go
package test

import (
	"github.com/homeport/pina-golada/pkg/files"
)

// TestInjectionPoint the variable in which the compiled interface implementation will be injected
var TestInjectionPoint TestInterface

// TestInterface is a testing interface for Pina-Golada. This interface will be implemented by Pina-Golada.
// The injector is the name of an exported variable of the type of this interface.
// The instance of the compiled struct, implementing this interface, will stored in the variable provided in the
// annotation below.
//
// @pgl(injector=TestInjectionPoint)
type TestInterface interface {
	
	// GetAssetFile is the method that returns a virtual asset stored in the directory instance.
	// The asset path passed to the annotation is relative to the location the tool will be called in.
	// The compressor defines what type of compression will be used to store the asset found as binary.
	// @pgl(asset=/assets/default-config.yaml&compressor=tar)
	GetAssetFile() (dir files.Directory , e error)
	
}
```

Now run `pina-golada generate` prior to your build. This will generate all needed files.
After you build your go application, use `pina-golada cleanup` to clean your project from the asset provider files.

## License

Licensed under [MIT License](https://github.com/homeport/pina-golada/blob/master/LICENSE)
