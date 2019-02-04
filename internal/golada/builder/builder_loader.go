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

package builder

import (
	"github.com/homeport/pina-golada/pkg/files"
	"github.com/homeport/pina-golada/pkg/files/paths"
	"os"
	"path/filepath"
)

// Load from disk loads the file into the directory from the disk
func loadFromDisk(root files.Directory, path string) (isDir bool, directory files.Directory, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, nil, err
	}

	if info.IsDir() {
		subDirectory := root.NewDirectory(paths.Of(filepath.Base(path))) // Load into sub path so we can store the name
		err := files.LoadFromDisk(subDirectory, path)
		if err != nil {
			return false, nil, err
		}
		return true, root, nil
	}

	if err := files.LoadFromDisk(root, path); err != nil {
		return false, nil, err
	}
	return false, root, nil
}
