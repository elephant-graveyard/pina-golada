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

package integration

import (
	"bytes"
	"github.com/homeport/pina-golada/internal/golada/builder"
	"github.com/homeport/pina-golada/internal/golada/cmd"
	"github.com/homeport/pina-golada/internal/golada/logger"
	"github.com/homeport/pina-golada/pkg/annotation"
	test "github.com/homeport/pina-golada/test/integration/.test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"os/exec"
	"testing"
)

func TestCompileIntegrationTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pgl integration compile")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	Expect(os.Chdir(".test")).To(Not(HaveOccurred())) // Change into the hidden test directory to target it

	Expect(os.Chmod("assets/file.txt" , 0733)).To(Not(HaveOccurred()))
	Expect(os.Chmod("assets/folder" , 0723)).To(Not(HaveOccurred()))
	Expect(os.Chmod("assets/folder/content.md" , 0722)).To(Not(HaveOccurred()))

	cmd.Generate(".", annotation.NewPropertyParser(), newLogger())

	Expect(os.Chdir("..")).To(Not(HaveOccurred())) // Return to the original directory
	return make([]byte, 0)
}, func(ignored []byte) {})

var _ = SynchronizedAfterSuite(func() {}, func() {
	cmd.Cleanup(".test/", newLogger())
})

var _ = Describe("should compile provided unit test", func() {

	var (
		buffer *bytes.Buffer
		writer test.BufferedWriter
	)

	_ = BeforeEach(func() {
		buffer = &bytes.Buffer{}
		writer = test.BufferedWriter{Buffer: buffer}
	})

	_ = It("should run the unit-tests in the compiled integration test", func() {
		Expect(os.Chdir(".test")).To(Not(HaveOccurred())) // Change into the hidden test directory to target it

		command := exec.Command("go", "test", "./...", "--", "count=1")
		command.Stdout = writer
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			_, err := os.Stdout.Write(buffer.Bytes())
			Expect(err).To(Not(HaveOccurred()))
			Fail("failed integration tests, check above")
			return
		}

		Expect(os.Chdir("..")).To(Not(HaveOccurred())) // Return to the original directory
	})
})

func newLogger() *logger.DefaultLogger {
	return logger.NewDefaultLogger(&builder.DevNullWriter{}, logger.Info)
}
