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
	"github.com/homeport/pina-golada/pkg/files/paths"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pgl .test sample")
}

var _ = AfterSuite(func() {
	Expect(os.RemoveAll("assets/test")).To(Not(HaveOccurred()))
})

var _ = Describe("Should have generated correctly", func() {
	_ = It("should compile file correctly", func() {
		dir, e := Provider.GetFileAsset()
		Expect(e).To(Not(HaveOccurred()))
		Expect(dir).To(Not(BeNil()))

		file := dir.File(paths.Of("file.txt"))
		Expect(file).To(Not(BeNil()))
		if NotWindows() {
			Expect(file.PermissionSet().Perm()).To(BeEquivalentTo(0733))
		}

		buffer := &bytes.Buffer{}
		Expect(file.CopyContent(buffer)).To(Not(HaveOccurred()))
		Expect(buffer.String()).To(BeEquivalentTo("All Your Base Are Belong to Us"))
	})

	_ = It("should have compiled folder correctly", func() {
		dir, e := Provider.GetFolderAsset()
		Expect(e).To(Not(HaveOccurred()))
		Expect(dir).To(Not(BeNil()))

		Expect(dir.Name().String()).To(BeEquivalentTo("folder"))

		content := dir.File(paths.Of("content.md"))
		Expect(dir).To(Not(BeNil()))

		buffer := &bytes.Buffer{}
		Expect(content.CopyContent(buffer)).To(Not(HaveOccurred()))
		Expect(buffer.String()).To(BeEquivalentTo("Hello there. General Kenobi."))
	})

	_ = It("should write files correctly", func() {
		dir, e := Provider.GetFileAsset()
		Expect(e).To(Not(HaveOccurred()))
		Expect(dir).To(Not(BeNil()))

		Expect(files.WriteToDisk(dir, "assets/test/", true)).To(Not(HaveOccurred()))

		file, e := ioutil.ReadFile("assets/test/file.txt")
		Expect(e).To(Not(HaveOccurred()))
		Expect(string(file)).To(BeEquivalentTo("All Your Base Are Belong to Us"))

		if !IsOS("windows") {
			Expect(GetFilePermission("assets/test/file.txt").Perm()).To(BeEquivalentTo(0733))
		}
	})

	_ = It("should write directories correctly", func() {
		dir, e := Provider.GetFolderAsset()
		Expect(e).To(Not(HaveOccurred()))
		Expect(dir).To(Not(BeNil()))

		Expect(files.WriteToDisk(dir, "assets/test/", true)).To(Not(HaveOccurred()))

		file, e := ioutil.ReadFile("assets/test/folder/content.md")
		Expect(e).To(Not(HaveOccurred()))
		Expect(string(file)).To(BeEquivalentTo("Hello there. General Kenobi."))

		if !IsOS("windows") {
			Expect(GetFilePermission("assets/test/folder/content.md").Perm()).To(BeEquivalentTo(0722))
			Expect(GetFilePermission("assets/test/folder").Perm()).To(BeEquivalentTo(0723))
		}
	})
})

// NotWindows returns if the system is not windows
func NotWindows() bool {
	return !IsOS("windows")
}
