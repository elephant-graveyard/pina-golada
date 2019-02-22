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

package compressor

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"path/filepath"

	"github.com/homeport/pina-golada/pkg/files"
	"github.com/homeport/pina-golada/pkg/files/paths"
)

// Tar is an implementation of the compressor interface which compresses to .tar.gz files
type Tar struct{}

// Compress compresses the directory into the writer
func (t *Tar) Compress(directory files.Directory, writer *bytes.Buffer) (e error) {
	gzipWriter := gzip.NewWriter(writer)
	tarWriter := tar.NewWriter(gzipWriter)
	files.WalkFileTree(directory, func(file files.File) {
		buffer := &bytes.Buffer{}
		if err := file.CopyContent(buffer); err != nil {
			return
		}

		tarHeader := &tar.Header{
			Name:     filepath.ToSlash(file.AbsolutePath().String()),
			Mode:     int64(file.PermissionSet()),
			Size:     int64(buffer.Len()),
			Typeflag: tar.TypeReg,
		}

		if err := tarWriter.WriteHeader(tarHeader); err != nil {
			return
		}
		if _, err := tarWriter.Write(buffer.Bytes()); err != nil {
			return
		}
	})

	if err := gzipWriter.Close(); err != nil {
		return err
	}
	return tarWriter.Close()
}

// Decompress decompresses the reader into the directory
func (t *Tar) Decompress(reader io.Reader) (directory files.Directory, e error) {
	root := files.NewRootDirectory()

	gzipReader, e := gzip.NewReader(reader)
	if e != nil {
		return root, e
	}

	tarReader := tar.NewReader(gzipReader)

	var foundError error
	for {
		header, bufferReaderError := tarReader.Next()
		if bufferReaderError == io.EOF {
			break
		}
		if bufferReaderError != nil {
			foundError = bufferReaderError
			break
		}

		if err := root.NewFile(paths.Of(header.Name)).WithPermission(header.FileInfo().Mode()).Write(tarReader); err != nil {
			foundError = err
			break
		}
	}

	if err := gzipReader.Close(); err != nil { // Close the gzip reader
		return nil, err
	}

	return root, foundError
}
