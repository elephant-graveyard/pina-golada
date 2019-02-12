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

package files

import (
	"io"
	"io/ioutil"

	"github.com/homeport/pina-golada/pkg/files/paths"
)

// File represents a digital file object
//
// Name returns the name of the file, excluding the paths it is stored in
//
// CopyContent stores the files binary content in the writer
//
// Write writes content to the file. The default will not append to the file
//
// Delete deletes the file
type File interface {
	Name() (name paths.Path)
	AbsolutePath() (path paths.Path)

	CopyContent(writer io.Writer) (e error)
	Write(reader io.Reader) (e error)
	WriteFlagged(reader io.Reader, append bool) (e error)
	Delete()

	Parent() (parentDirectory Directory)
}

// memoryFile is an in memory implementation of the File interface
type memoryFile struct {
	name    paths.Path
	parent  Directory
	content []byte
}

// Name Returns the name of the file
func (m *memoryFile) Name() (name paths.Path) {
	return m.name
}

// Returns the absolute path
func (m *memoryFile) AbsolutePath() (path paths.Path) {
	return m.Parent().AbsolutePath().Concat(m.Name())
}

// CopyContent copies the content of the writer
func (m *memoryFile) CopyContent(writer io.Writer) (e error) {
	_, e = writer.Write(m.content)
	return e
}

// Write writes the content of the reader
func (m *memoryFile) Write(reader io.Reader) (e error) {
	return m.WriteFlagged(reader, false)
}

// WriteFlagged writes the content of the reader to the file and appends it if appendBytes is true
func (m *memoryFile) WriteFlagged(reader io.Reader, appendBytes bool) (e error) {
	bytes, e := ioutil.ReadAll(reader)
	if e != nil {
		return e
	}

	if appendBytes {
		m.content = append(m.content, bytes...)
	} else {
		m.content = bytes
	}
	return nil
}

// Delete deletes the file
func (m *memoryFile) Delete() {
	m.Parent().DeleteFile(m.Name())
}

// Parent returns the parent of the file
func (m *memoryFile) Parent() (parentDirectory Directory) {
	return m.parent
}
