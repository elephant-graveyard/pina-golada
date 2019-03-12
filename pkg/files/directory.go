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
	"bytes"
	"github.com/homeport/pina-golada/pkg/files/paths"
	"os"
)

// Directory represents a virtual directory containing files.
//
// Name return the name of the directory excluding the paths it is in
//
// Files returns the list of all files in the given directory
//
// File returns the file for the given name or nil
//
// NewFile creates a new File in the directory and returns the new instance
//
// Delete deletes the file stored under the given name
//
// Directories returns the directories stored under this directory
//
// Directory returns the directory with the given name or nil
//
// NewDirectory creates a new directory
//
// DeleteDirectory deletes a directory
//
// Parent returns the directory this directory is found in
//
// AsRoot creates a deep copy of the current directory, but with the current directory as it's root
type Directory interface {
	Name() (name paths.Path)
	AbsolutePath() (path paths.Path)

	WithPermission(permission os.FileMode) Directory
	PermissionSet() os.FileMode

	Files() (files []File)
	File(path paths.Path) (file File)
	NewFile(path paths.Path) (newFile File)
	DeleteFile(path paths.Path)

	Directories() (directories []Directory)
	Directory(path paths.Path) (directory Directory)
	NewDirectory(path paths.Path) (newDirectory Directory)
	DeleteDirectory(path paths.Path)

	Parent() (parentDirectory Directory)
	AsRoot() (rootDirectory Directory)
}

// memoryDirectory is a in memory implementation of the directory interface
type memoryDirectory struct {
	name     paths.Path
	parent   Directory
	files    []File
	dirs     []Directory
	PermBits os.FileMode
}

// Name returns the name of the directory
func (m *memoryDirectory) Name() (name paths.Path) {
	return m.name
}

// Returns the absolute path of the directory
func (m *memoryDirectory) AbsolutePath() (path paths.Path) {
	if m.Parent() != nil {
		return m.Parent().AbsolutePath().Concat(m.Name())
	}
	return m.Name()
}

// WithPermission stores the permission set on the directory
func (m *memoryDirectory) WithPermission(permission os.FileMode) Directory {
	m.PermBits = permission
	return m
}

// PermissionSet returns the permission set of the directory
func (m *memoryDirectory) PermissionSet() os.FileMode {
	return m.PermBits
}

// Files returns a slice of all files found in the directory
func (m *memoryDirectory) Files() (files []File) {
	return m.files
}

// File returns the file for the given name
func (m *memoryDirectory) File(path paths.Path) (file File) {
	if !path.Valid() {
		return nil
	}

	if path.Direct() {
		for _, file := range m.Files() {
			if file.Name().Equals(path) {
				return file
			}
		}
	} else {
		newPath := path.Clone()

		subDirectory := m.Directory(newPath.Pop())
		if subDirectory != nil {
			return subDirectory.File(newPath)
		}
	}

	return nil
}

// NewFile creates a new file at the given path
func (m *memoryDirectory) NewFile(path paths.Path) (newFile File) {
	if !path.Valid() {
		return nil
	}

	newPath := path.Clone()
	if !path.Direct() {
		fileName := newPath.Drop()
		return m.NewDirectory(newPath).NewFile(fileName)
	}

	if foundFile := m.File(newPath); foundFile != nil { // return existing file
		return foundFile
	}

	file := &memoryFile{
		parent:   m,
		name:     newPath,
	}
	file.WithPermission(m.PermissionSet())
	m.files = append(m.files, file)
	return file
}

// DeleteFile deletes a file from the directory
func (m *memoryDirectory) DeleteFile(path paths.Path) {
	if !path.Valid() {
		return
	}

	if !path.Direct() {
		newPath := path.Clone()
		fileName := newPath.Drop()
		parentDirectory := m.Directory(newPath)
		if parentDirectory != nil {
			parentDirectory.DeleteFile(fileName)
		}
		return
	}

	for index, file := range m.Files() {
		if file.Name().Equals(path) {
			m.files = append(m.files[:index], m.files[index+1:]...)
		}
	}
	return
}

// Directories lists all the directories in the directory
func (m *memoryDirectory) Directories() (directories []Directory) {
	return m.dirs
}

// Directory returns the directory
func (m *memoryDirectory) Directory(path paths.Path) (directory Directory) {
	if !path.Valid() {
		return nil
	}

	newPath := path.Clone()
	if !path.Direct() {
		firstDirectory := m.Directory(newPath.Pop())
		if firstDirectory == nil {
			return nil
		}

		return firstDirectory.Directory(newPath)
	}

	for _, dir := range m.Directories() {
		if dir.Name().Equals(path) {
			return dir
		}
	}

	return nil
}

// NewDirectory creates a new directory
func (m *memoryDirectory) NewDirectory(path paths.Path) (newDirectory Directory) {
	if !path.Valid() {
		return nil
	}

	if !path.Direct() {
		newPath := path.Clone()
		thisLevelDirectory := newPath.Pop()

		foundDirectory := m.Directory(thisLevelDirectory)
		if foundDirectory != nil {
			return foundDirectory.NewDirectory(newPath)
		}

		newLevelDirectory := m.NewDirectory(thisLevelDirectory)
		if newLevelDirectory != nil {
			return newLevelDirectory.NewDirectory(newPath)
		}
	}

	if foundDirectory := m.Directory(path); foundDirectory != nil {
		return foundDirectory
	}

	createdDirectory := &memoryDirectory{
		name:   path,
		parent: m,
	}
	createdDirectory.WithPermission(m.PermissionSet())

	m.dirs = append(m.dirs, createdDirectory)
	return createdDirectory
}

// DeleteDirectory deletes a directory
func (m *memoryDirectory) DeleteDirectory(path paths.Path) {
	if !path.Valid() {
		return
	}

	if !path.Direct() {
		newPath := path.Clone()
		directoryName := newPath.Drop()
		parentDirectory := m.Directory(newPath)
		if parentDirectory != nil {
			parentDirectory.DeleteDirectory(directoryName)
		}
		return
	}

	for index, directory := range m.Directories() {
		if directory.Name().Equals(path) {
			m.dirs = append(m.dirs[:index], m.dirs[index+1:]...)
		}
	}
	return
}

// Parent returns the parent directory
func (m *memoryDirectory) Parent() (parentDirectory Directory) {
	return m.parent
}

// AsRoot creates a deep copy of the current directory, but with the current directory as it's root
func (m *memoryDirectory) AsRoot() (rootDirectory Directory) {
	root := NewRootDirectory()
	if err := copyDirectory(m, root); err != nil {
		return nil
	}
	return root
}

// copyDirectory copies the content of one directory into the other
func copyDirectory(original Directory, new Directory) error {
	for _, f := range original.Files() {
		content := &bytes.Buffer{}
		if err := f.CopyContent(content); err != nil {
			return err
		}

		if err := new.NewFile(f.Name()).Write(content); err != nil {
			return nil
		}
	}

	for _, dir := range original.Directories() {
		if err := copyDirectory(dir, new.NewDirectory(dir.Name())); err != nil {
			return err
		}
	}
	return nil
}

// NewRootDirectory returns a new root directory
func NewRootDirectory() Directory {
	return &memoryDirectory{
		name:     paths.RootPath(),
	}
}
