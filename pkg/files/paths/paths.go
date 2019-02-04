// Copyright © 2019 The Homeport Team

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package paths

import (
	"os"
	"path/filepath"
	"strings"
)

// Path represents a relative paths to a file
//
// ElementAt returns the element of the paths at the given index starting at 0
// If index is out of range, the returned path is not valid
//
// Direct returns if the paths only has one element left
//
// Valid returns if the path is actually valid
//
// Size returns the amount of elements left in the paths
//
// Clone will create new instance of the paths
//
// Concat will concat the two paths together
//
// Slice returns the slice of path entries
//
// Pop will pop the first element in the paths
//
// Drop drops the last element in the path
type Path interface {
	ElementAt(index int) Path
	Direct() bool
	Valid() bool
	Size() int

	Clone() Path
	Equals(path Path) bool

	Concat(other Path) Path
	Slice() []string
	String() string

	Pop() Path
	Drop() Path
}

// MemoryPath implements the Path interface
type MemoryPath struct {
	path []string
}

// ElementAt returns the element of the paths at the given index starting at 0
// If index is out of range, return an empty string
func (m *MemoryPath) ElementAt(index int) Path {
	if index < 0 || index >= m.Size() {
		return &MemoryPath{}
	}
	return uncheckedOfString(m.path[index])
}

// Direct returns if the paths only has one element left
func (m *MemoryPath) Direct() bool {
	return m.Size() == 1
}

// Valid returns if the path is actually valid
func (m *MemoryPath) Valid() bool {
	return m.Size() > 0
}

// Size returns the size of the paths
func (m *MemoryPath) Size() int {
	return len(m.path)
}

// Clone will create new instance of the paths
func (m *MemoryPath) Clone() Path {
	return &MemoryPath{
		path: m.path,
	}
}

// Equals returns if the path is equal to another
func (m *MemoryPath) Equals(path Path) bool {
	if m.Size() != path.Size() {
		return false
	}

	otherSlice := path.Slice()
	for index, pathEntry := range m.Slice() {
		if otherSlice[index] != pathEntry {
			return false
		}
	}

	return true
}

// Concat will concat the two paths together
func (m *MemoryPath) Concat(other Path) Path {
	return &MemoryPath{
		path: append(m.path, other.Slice()...),
	}
}

// Slice returns the slice of path entries
func (m *MemoryPath) Slice() []string {
	return m.path
}

// Pop will pop the first element in the paths
func (m *MemoryPath) Pop() Path {
	if m.Size() > 0 {
		firstElement := m.path[0]
		m.path = m.path[1:]
		return uncheckedOfString(firstElement)
	}
	return &MemoryPath{}
}

// Drop will drop the last element in the path
func (m *MemoryPath) Drop() Path {
	if m.Size() > 0 {
		lastElement := m.path[m.Size()-1]
		m.path = m.path[:m.Size()-1]
		return uncheckedOfString(lastElement)
	}
	return &MemoryPath{}
}

// Returns the paths as a string
func (m *MemoryPath) String() string {
	return strings.Join(m.path, GetPathSeparator())
}

// Of creates a paths off of a string
func Of(path string) *MemoryPath {
	cleansedPath := strings.Trim(filepath.Clean(filepath.FromSlash(path)), GetPathSeparator())
	return &MemoryPath{
		path: strings.Split(cleansedPath, GetPathSeparator()),
	}
}

// uncheckedOfString creates a path with one entry
func uncheckedOfString(path string) *MemoryPath {
	return &MemoryPath{
		path: []string{path},
	}
}

// GetPathSeparator returns the paths separator of the current os
func GetPathSeparator() string {
	return string(os.PathSeparator)
}

// RootPath returns a path with one single empty string entry
func RootPath() Path {
	return uncheckedOfString("")
}
