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

package files

import (
	"bytes"
	"github.com/homeport/pina-golada/pkg/files/paths"
	"io/ioutil"
	"os"
	"path/filepath"
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

// LoadFromDisk Loads the content of the paths into the directory recursively
func LoadFromDisk(directory Directory, path string) (e error) {
	info, statError := os.Stat(path)
	if statError != nil {
		return statError
	}

	if info.IsDir() {
		directoryContent, e := ioutil.ReadDir(path)
		if e != nil {
			return e
		}

		for _, file := range directoryContent {
			if file.IsDir() {
				if err := LoadFromDisk(directory.NewDirectory(paths.Of(file.Name())), filepath.Join(path, file.Name())); err != nil {
					return err
				}
			} else {
				if err := readFileInto(directory, filepath.Join(path, file.Name())); err != nil {
					return err
				}
			}
		}
	} else {
		if err := readFileInto(directory, path); err != nil {
			return err
		}
	}
	return nil
}

func readFileInto(directory Directory, path string) (e error) {
	content, readError := ioutil.ReadFile(path)
	if readError != nil {
		return readError
	}

	return directory.NewFile(paths.Of(path).Drop()).Write(bytes.NewBuffer(content))
}
