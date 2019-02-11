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

package inspector

import (
	"os"
	"path/filepath"
)

// FileStream defines a stream of files based on a provided file supplier
type FileStream struct {
	Files []File
}

// Count returns the count
func (f *FileStream) Count() (count int) {
	return len(f.Files)
}

// Filter filters the file stream and removes all files that don't match the filter
func (f *FileStream) Filter(filter func(file File) bool) *FileStream {
	files := make([]File, 0)
	for _, f := range f.Files {
		if filter(f) {
			files = append(files, f)
		}
	}
	f.Files = files
	return f
}

// ForEach executes the consumer for each file in the stream
func (f *FileStream) ForEach(consumer func(file File)) {
	for _, file := range f.Files {
		consumer(file)
	}
}

// File is a struct containing the reference to a file as well as it's info
type File struct {
	FileInfo os.FileInfo
	Path     string
}

// NewFileStream creates a new file stream based on the provided paths
func NewFileStream(path string) (stream *FileStream, e error) {
	files := make([]File, 0)
	er := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		files = append(files, File{
			FileInfo: info,
			Path:     path,
		})
		return err
	})
	return &FileStream{Files: files}, er
}
