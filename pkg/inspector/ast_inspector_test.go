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

package inspector

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

// Test han shot first
var Test = ""

// TestInterface black magic
// its black magic i swear
type TestInterface interface {
	// Foo test
	Foo() bool
}

var _ = Describe("Testing default ast stream functionality", func() {
	_ = It("should find the test interface defined here", func() {
		stream, e := NewFileStream("./")
		Expect(e).To(BeNil())

		astStream := NewAstStream(stream.Filter(func(file File) bool {
			return strings.Contains(file.FileInfo.Name(), "ast_inspector_test.go")
		}))
		interfaces := astStream.Find()
		Expect(len(interfaces)).To(BeEquivalentTo(1))

		foundInterface := interfaces[0]
		Expect(foundInterface.Name.Name).To(BeEquivalentTo("TestInterface"))
		Expect(foundInterface.Docs.Text()).To(ContainSubstring("TestInterface black magic"))

		methods := foundInterface.InterfaceReference.Methods
		Expect(methods.NumFields()).To(BeEquivalentTo(1))
		Expect(methods.List[0].Doc.Text()).To(ContainSubstring("Foo test"))

	})
})
