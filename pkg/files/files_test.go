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
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/homeport/pina-golada/pkg/files/paths"
)

func TestFiles(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pgl pkg files")
}

var _ = Describe("should handle files properly", func() {

	var (
		root           Directory
		buffer         *bytes.Buffer
		testGoFile     paths.Path
		testHomeFolder paths.Path
	)

	BeforeEach(func() {
		root = NewRootDirectory()
		buffer = &bytes.Buffer{}
		testGoFile = paths.Of("/usr/homeport/home/test.go")
		testHomeFolder = paths.Of("/usr/homeport/home")
	})

	It("should create memory files correctly", func() {
		e := root.NewFile(paths.Of("test.txt")).Write(bytes.NewBufferString("test"))
		Expect(e).To(BeNil())
		Expect(len(root.Files())).To(BeEquivalentTo(1))

		file := root.File(paths.Of("test.txt"))
		Expect(file).To(Not(BeNil()))

		Expect(file.CopyContent(buffer)).To(BeNil())
		Expect(string(buffer.Bytes())).To(BeEquivalentTo("test"))
	})

	_ = It("should return the correct paths", func() {
		testFile := root.NewDirectory(paths.Of("usr")).NewDirectory(paths.Of("homeport")).NewDirectory(paths.Of("home")).NewFile(paths.Of("test.go"))
		Expect(testFile.AbsolutePath().String()).To(BeEquivalentTo(filepath.FromSlash("/usr/homeport/home/test.go")))
	})

	_ = It("should create nested files", func() {

		e := root.NewFile(testGoFile).Write(bytes.NewBufferString("test"))
		Expect(e).To(BeNil())

		file := root.File(testGoFile)
		Expect(file).To(Not(BeNil()))
		Expect(file.CopyContent(buffer)).To(BeNil())
		Expect(buffer.String()).To(BeEquivalentTo("test"))
	})

	_ = It("should not find that file", func() {
		Expect(root.File(paths.Of("usr/whatever/this/is/test.ruby"))).To(BeNil())
	})

	_ = It("should delete filed properly", func() {
		_ = root.NewFile(testGoFile).Write(bytes.NewBufferString("test"))

		f := root.File(testGoFile)
		Expect(f).To(Not(BeNil()))

		root.DeleteFile(testGoFile)
		f = root.File(testGoFile)
		Expect(f).To(BeNil())
	})

	_ = It("should delete filed properly from files", func() {
		_ = root.NewFile(testGoFile).Write(bytes.NewBufferString("test"))

		f := root.File(testGoFile)
		Expect(f).To(Not(BeNil()))

		f.Delete()
		f = root.File(testGoFile)
		Expect(f).To(BeNil())
	})

	_ = It("should delete directories properly", func() {
		_ = root.NewDirectory(testHomeFolder)

		d := root.Directory(testHomeFolder)
		Expect(d).To(Not(BeNil()))

		root.DeleteDirectory(testHomeFolder)
		d = root.Directory(testHomeFolder)
		Expect(d).To(BeNil())
	})

	_ = It("should delete files with a complex paths", func() {
		root.NewFile(testGoFile)

		directory := root.Directory(testHomeFolder)
		Expect(directory).To(Not(BeNil()))
		Expect(len(directory.Files())).To(BeEquivalentTo(1))

		root.DeleteFile(testGoFile)
		Expect(len(directory.Files())).To(BeEquivalentTo(0))
	})

	_ = It("should not create multiple files", func() {
		root.NewFile(testGoFile)
		root.NewFile(testGoFile)

		Expect(len(root.Directories())).To(BeEquivalentTo(1))

		homeDirectory := root.Directory(testHomeFolder)
		Expect(homeDirectory).To(Not(BeNil()))
		Expect(len(homeDirectory.Files())).To(BeEquivalentTo(1))
	})

	_ = It("should append properly", func() {
		newFile := root.NewFile(testGoFile)
		Expect(newFile.Write(bytes.NewBufferString("te"))).To(BeNil())
		Expect(newFile.WriteFlagged(bytes.NewBufferString("st"), true)).To(BeNil())

		Expect(newFile.CopyContent(buffer)).To(BeNil())

		Expect(buffer.String()).To(BeEquivalentTo("test"))
	})

	_ = It("should create an asRoot copy", func() {
		root.NewDirectory(paths.Of("usr")).NewDirectory(paths.Of("homeport")).NewDirectory(paths.Of("home")).NewFile(paths.Of("test.go"))
		rootCopy := root.Directory(paths.Of("usr/homeport")).AsRoot()
		Expect(rootCopy.File(paths.Of("home/test.go"))).To(Not(BeNil()))
	})

	_ = It("should read from file correctly", func() {
		e := LoadFromDisk(root, filepath.FromSlash("../../assets/tests"))
		Expect(e).To(BeNil())

		fileTestFolder := root.Directory(paths.Of("fileTestFolder"))
		Expect(fileTestFolder).To(Not(BeNil()))

		file := fileTestFolder.File(paths.Of("test.txt"))
		Expect(file).To(Not(BeNil()))
		Expect(file.CopyContent(buffer)).To(BeNil())

		Expect(strings.Replace(buffer.String(), "\r", "", -1)).To(BeEquivalentTo("test\n")) // remove \r from windows
	})

	_ = It("should overwrite permission sets on files", func() {
		pathToTestDir := "../../assets/tests/issue-51/file"
		pathToTestFile := filepath.Join(pathToTestDir, "permission.txt")

		os.Chmod(pathToTestFile, 0777)
		Expect(GetFilePermission(pathToTestFile)).To(BeEquivalentTo(0777))

		Expect(LoadFromDisk(root, pathToTestDir)).To(BeNil())
		file := root.File(paths.Of("permission.txt"))
		Expect(file).To(Not(BeNil()))

		file.WithPermission(os.FileMode(0701))
		Expect(WriteToDisk(root, pathToTestDir, true)).To(BeNil())
		Expect(GetFilePermission(pathToTestFile)).To(BeEquivalentTo(os.FileMode(0701)))
	})

	_ = It("should overwrite permission sets on folders", func() {
		pathToTestDir := "../../assets/tests/issue-51/folder"
		pathToTestFile := filepath.Join(pathToTestDir, "permission")

		Expect(os.Chmod(pathToTestFile, 0777)).To(Not(HaveOccurred()))
		Expect(GetFilePermission(pathToTestFile).Perm()).To(BeEquivalentTo(0777))

		Expect(LoadFromDisk(root, pathToTestDir)).To(BeNil())
		directory := root.Directory(paths.Of("permission"))
		Expect(directory).To(Not(BeNil()))

		directory.WithPermission(0701)
		Expect(WriteToDisk(root, pathToTestDir, true)).To(BeNil())
		Expect(GetFilePermission(pathToTestFile).Perm()).To(BeEquivalentTo(0701))
	})
})

func GetFilePermission(path string) os.FileMode {
	info, e := os.Stat(path)
	Expect(e).To(Not(HaveOccurred()))
	return info.Mode()
}
