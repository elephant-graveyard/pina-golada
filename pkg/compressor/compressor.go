// Copyright © 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package compressor hosts the compressor logic used by pina-golada to compress the assets
// into binary form. In here, the default compressors pina-golada ships with are implemented.
package compressor

import (
	"bytes"
	"io"
	"strings"

	"github.com/homeport/pina-golada/pkg/files"
)

var (
	// DefaultRegistry contains the registered compressors
	DefaultRegistry = NewMapRegistry().Put("tar", &Tar{})
)

// Registry contains a collection of Compressor instances that can be used to compress files
//
// Put stores a new compressor instance in the registry for the given id
//
// Find returns the compressor that is registered for the given id or nil if non was found
type Registry interface {
	Put(id string, compressor Compressor) Registry
	Find(id string) (compressor Compressor)
}

// MapRegistry is a map based compressor registry
type MapRegistry struct {
	tracking map[string]Compressor
}

// NewMapRegistry creates a new instance of the map registry
func NewMapRegistry() *MapRegistry {
	return &MapRegistry{
		tracking: make(map[string]Compressor),
	}
}

// Put stores a new compressor instance for the given id in the registry
func (r *MapRegistry) Put(id string, compressor Compressor) Registry {
	r.tracking[strings.ToLower(id)] = compressor
	return r
}

// Find returns the found compressor instance for the given id, or nil if non was found
func (r *MapRegistry) Find(id string) (compressor Compressor) {
	return r.tracking[strings.ToLower(id)]
}

// Compressor defines an object that is capable of compressing and decompressing a file
type Compressor interface {
	// Compress compresses the provided filepath into the output byte slice

	Compress(directory files.Directory, writer *bytes.Buffer) (e error)

	// Decompress decompresses the file paths into the output paths
	Decompress(reader io.Reader) (dir files.Directory, e error)
}
