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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/homeport/pina-golada/pkg/files/paths"
)

// WalkFileTree iterates over each and every file instance found in the directory
func WalkFileTree(directory Directory, consumer func(file File)) {
	for _, file := range directory.Files() {
		consumer(file)
	}

	for _, dir := range directory.Directories() {
		WalkFileTree(dir, consumer)
	}
}

// WalkDirectoryTree walks over each sub directory found in the directory, excluding the passed directory
func WalkDirectoryTree(directory Directory, consumer func(d Directory)) {
	for _, dir := range directory.Directories() {
		consumer(dir)
		WalkDirectoryTree(dir, consumer)
	}
}

// LoadFromDisk loads the content of the paths into the directory recursively
func LoadFromDisk(directory Directory, path string) (e error) {
	_, e = LoadFromDiskAndType(directory, path)
	return e
}

// LoadFromDiskAndType loads the content of the paths into the directory recursively.
// This also returns true if the file located under the path is a directory or false
// if the file at the path is a flat file.
func LoadFromDiskAndType(directory Directory, path string) (isDir bool, e error) {
	info, statError := os.Stat(path)
	if statError != nil {
		return false, statError
	}

	if info.IsDir() {
		directory.WithPermission(info.Mode())

		directoryContent, e := ioutil.ReadDir(path)
		if e != nil {
			return true, e
		}

		for _, file := range directoryContent {
			if file.IsDir() {
				if err := LoadFromDisk(directory.NewDirectory(paths.Of(file.Name())).WithPermission(file.Mode()), filepath.Join(path, file.Name())); err != nil {
					return true, err
				}
			} else {
				if err := readFileInto(directory, filepath.Join(path, file.Name())); err != nil {
					return true, err
				}
			}
		}
		return true, nil
	}

	if err := readFileInto(directory, path); err != nil {
		return false, err
	}
	return false, nil
}

func readFileInto(directory Directory, path string) (e error) {
	content, readError := ioutil.ReadFile(path)
	if readError != nil {
		return readError
	}

	info, statError := os.Stat(path)
	if statError != nil {
		return statError
	}

	return directory.NewFile(paths.Of(path).Drop()).WithPermission(info.Mode()).Write(bytes.NewBuffer(content))
}

// WriteToDisk writes a directory to the given path
// Overwriting or skipping an existing file based on the bool
func WriteToDisk(directory Directory, path string, overwrite bool) (e error) {
	info, statError := os.Stat(path)
	if statError != nil {
		if !os.IsNotExist(statError) {
			return statError
		}

		if err := os.MkdirAll(path, directory.PermissionSet()); err != nil {
			return err
		}
	}

	if info != nil && !info.IsDir() {
		return fmt.Errorf("provided path pointed to file %s", path)
	}

	if directory.Parent() != nil { // If the directory is not a root directory, we want to create the directory
		path = filepath.Join(path, directory.Name().String())
	}

	return writeDirectoryToDisk(directory, path, overwrite)
}

func writeFileToDisk(file File, directoryPath string, overwrite bool) (e error) {
	path := filepath.Join(directoryPath, file.Name().String())
	var fileOnDisk *os.File

	info, e := os.Stat(path)
	if e != nil { // Checking if statistics for path cannot be found
		if !os.IsNotExist(e) { // Checking if the path does exist but couldn't be accessed
			return e
		}

		if fileOnDisk, e = os.Create(path); e == nil { // Creating the path given that it doesn't exist
			if err := fileOnDisk.Close(); err != nil {
				return err
			}
		} else {
			return e
		}
	}

	if info != nil { // Checking if the path currently exists
		if info.IsDir() { // Checking if the path is a directory
			return fmt.Errorf("provided path pointed to directory %s", path)
		}

		if !overwrite { // Checking if the file should be overwritten given it exists
			return nil
		}
	}

	if fileOnDisk, e = os.OpenFile(path, os.O_RDWR|os.O_TRUNC, file.PermissionSet()); e != nil {
		// Check for errors after opening the file
		return e
	}

	if err := file.CopyContent(fileOnDisk); err != nil {
		// Check for errors after writing the given file to the disk file
		return err
	}

	if err := fileOnDisk.Close(); err != nil {
		// Checking for errors while closing the disk file
		return err
	}
	return nil
}

func writeDirectoryToDisk(directory Directory, directoryPath string, overwrite bool) (e error) {
	info, e := os.Stat(directoryPath)
	if e != nil {
		if !os.IsNotExist(e) {
			return e
		}

		if e := os.MkdirAll(directoryPath, directory.PermissionSet()); e != nil {
			return e
		}
	}

	if info != nil && !info.IsDir() {
		return fmt.Errorf("provided path pointed to file %s", directoryPath)
	}

	for _, dir := range directory.Directories() {
		if err := writeDirectoryToDisk(dir, filepath.Join(directoryPath, dir.Name().String()), overwrite); err != nil {
			return err
		}
	}

	for _, file := range directory.Files() {
		if err := writeFileToDisk(file, directoryPath, overwrite); err != nil {
			return err
		}
	}

	return nil
}
