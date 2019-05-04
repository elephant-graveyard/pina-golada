/*
 * Copyright © 2019 The Homeport Team
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package test

import (
	"bytes"
	"github.com/homeport/pina-golada/pkg/files"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"runtime"
)

var Provider Assets

// Assets is the test injector variable
// @pgl(injector=Provider)
type Assets interface {
	// @pgl(asset=assets/file.txt&compressor=tar)
	GetFileAsset() (dir files.Directory, e error)

	// @pgl(asset=assets/folder&compressor=tar)
	GetFolderAsset() (dir files.Directory, e error)
}

// IsOS returns if the current os equals the string
func IsOS(os string) bool {
	return runtime.GOOS == os
}

// DevNull defines a logger that just ignores all output
type DevNull struct{}

// Write writes the bytes in the DevNull buffer
func (*DevNull) Write(p []byte) (n int, err error) {
	return 0, nil
}

// GetFilePermission returns the permission of the file
func GetFilePermission(path string) os.FileMode {
	info, e := os.Stat(filepath.FromSlash(path))
	Expect(e).To(Not(HaveOccurred()))
	return info.Mode()
}

type BufferedWriter struct {
	Buffer *bytes.Buffer
}

func (b BufferedWriter) Write(p []byte) (n int, err error) {
	b.Buffer.Write(p)
	return len(p), nil
}
