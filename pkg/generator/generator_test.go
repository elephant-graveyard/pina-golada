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

package generator

import (
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pgl pkg generator")
}

var _ = Describe("should generate file correctly", func() {

	var (
		writer    *strings.Builder
		generator GoGenerator
	)

	_ = BeforeEach(func() {
		writer = &strings.Builder{}
		generator = NewFileGoGenerator("test_package")
	})

	_ = It("should generate an empty file correctly", func() {
		generator.Flush(writer)
		Expect(writer.String()).To(BeEquivalentTo("package test_package" + EndOfLine + EndOfLine))
	})

	_ = It("should import correctly", func() {
		generator.Import("strings").Import("os").Flush(writer)
		Expect(writer.String()).To(BeEquivalentTo(
			`package test_package` + EndOfLine + EndOfLine + "import (" + EndOfLine +
				Intend + Quote + "strings" + Quote + EndOfLine +
				Intend + Quote + "os" + Quote + EndOfLine +
				")" + EndOfLine + EndOfLine))
	})

	_ = It("should generate structs correctly", func() {
		generator.Struct("Test", func(s StructGenerator) {
			s.Field("TestBool", "bool")
		}).Flush(writer)

		Expect(writer.String()).To(BeEquivalentTo(
			`package test_package` + EndOfLine + EndOfLine + "type Test struct {" + EndOfLine +
				Intend + "TestBool bool" + EndOfLine +
				"}" + EndOfLine + EndOfLine))
	})

	_ = It("should generate methods correctly", func() {
		generator.Method("Add", func(m MethodGenerator) {
			m.Parameters("a int", "b int").ReturnTypes("c int").Body([]string{
				"return a + b",
			}...)
		}).Flush(writer)

		Expect(writer.String()).To(BeEquivalentTo(`package test_package` + EndOfLine + EndOfLine +
			"func Add(a int, b int) (c int) {" + EndOfLine +
			Intend + "return a + b" + EndOfLine +
			"}" + EndOfLine + EndOfLine))
	})

	_ = It("should generate receiver methods correctly", func() {
		generator.Method("Add", func(m MethodGenerator) {
			m.Receiver("s SuperCoolObject").Parameters("a int", "b int").ReturnTypes("c int").Body([]string{
				"return a + b",
			}...)
		}).Flush(writer)

		Expect(writer.String()).To(BeEquivalentTo(`package test_package` + EndOfLine + EndOfLine +
			"func (s SuperCoolObject) Add(a int, b int) (c int) {" + EndOfLine +
			Intend + "return a + b" + EndOfLine +
			"}" + EndOfLine + EndOfLine))
	})
})
