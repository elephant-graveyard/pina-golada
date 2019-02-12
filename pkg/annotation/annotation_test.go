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

package annotation

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAnnotations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pgl pkg annotation")
}

// TestAnnotation is just a really basic annotation
type TestAnnotation struct {
	SuperCoolName string `yaml:"name"`
	Version       int    `yaml:"version"`
}

// Returns the test annotation identifier
func (TestAnnotation) GetIdentifier() string {
	return "@test"
}

var _ = Describe("should parse annotations correctly", func() {

	var (
		parser     Parser
		annotation *TestAnnotation
	)

	_ = BeforeEach(func() {
		parser = NewCsvParser()
		annotation = &TestAnnotation{}
	})

	_ = It("should parse a normal annotation", func() {
		if err := parser.Parse("@test(name,test;version,1)", annotation); err != nil {
			log.Fatal(err)
			Fail("failed due to " + err.Error())
		}

		Expect(annotation.SuperCoolName).To(BeEquivalentTo("test"))
		Expect(annotation.Version).To(BeEquivalentTo(1))
	})

	_ = It("should load it from a comment string", func() {
		if err := parser.Parse(`MethodA does something this 
is a doc
@test(name,test;version,2)
More documentation`, annotation); err != nil {
			log.Fatal(err)
			Fail("failed due to " + err.Error())
		}

		Expect(annotation.SuperCoolName).To(BeEquivalentTo("test"))
		Expect(annotation.Version).To(BeEquivalentTo(2))
	})

	_ = It("should not find an annotation", func() {
		err := parser.Parse("This is just documentation @CoolAnnotationHe(a,g;1,e)", annotation)
		Expect(err).To(BeEquivalentTo(ErrNoAnnotation))
	})
})
