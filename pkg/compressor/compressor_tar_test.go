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

package compressor

import (
	"bytes"
	"github.com/homeport/pina-golada/pkg/files"
	"github.com/homeport/pina-golada/pkg/files/paths"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCompressor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pgl pkg compressor")
}

var _ = Describe("should compress files correctly", func() {
	var (
		registry  Registry
		directory files.Directory
		buffer    *bytes.Buffer
	)

	_ = BeforeEach(func() {
		registry = NewMapRegistry()
		registry.Put("tar", &Tar{})
		directory = files.NewRootDirectory()
		buffer = &bytes.Buffer{}
	})

	_ = It("should compress and decompress tars correctly", func() {
		Expect(directory.NewFile(paths.Of("usr/homeport/home/testA.go")).Write(bytes.NewBufferString("testA"))).To(BeNil())
		Expect(directory.NewFile(paths.Of("usr/homeport/home/testB.go")).Write(bytes.NewBufferString("testB"))).To(BeNil())

		tarCompressor := registry.Find("tar")
		Expect(tarCompressor).To(Not(BeNil()))

		Expect(tarCompressor.Compress(directory, buffer)).To(BeNil())

		result, e := tarCompressor.Decompress(buffer)
		Expect(e).To(BeNil())
		Expect(result).To(Not(BeNil()))

		testA := directory.File(paths.Of("usr/homeport/home/testA.go"))
		Expect(testA).To(Not(BeNil()))
		Expect(testA.CopyContent(buffer)).To(BeNil())
		Expect(buffer.String()).To(BeEquivalentTo("testA"))

		buffer.Reset()

		testB := directory.File(paths.Of("usr/homeport/home/testB.go"))
		Expect(testB).To(Not(BeNil()))
		Expect(testB.CopyContent(buffer)).To(BeNil())
		Expect(buffer.String()).To(BeEquivalentTo("testB"))
	})
})
