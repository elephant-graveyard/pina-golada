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

package cmd

import (
	"bufio"
	"github.com/homeport/pina-golada/internal/golada/builder"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var cleanupCommand = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleans all the generated files",
	Long:  "Cleanup will iterate over every go file in the project and delete it if it does start with PinaGolada's identifier string",
	Run: func(c *cobra.Command, args []string) {
		e := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if GoFileSelector.Match([]byte(info.Name())) {
				file, err := os.Open(path)
				if err != nil {
					return err
				}

				fileReader := bufio.NewReader(file)
				line, _, err := fileReader.ReadLine()
				if err != nil {
					_ = file.Close()
					return err
				}

				if err := file.Close(); err != nil {
					return err
				}
				if strings.EqualFold(string(line), builder.IdentifierString) {
					if err := os.Remove(path); err != nil {
						return err
					}
				}
			}
			return nil
		})
		if e != nil {
			log.Fatalf("could not iterrate over the files in this directory %s", e.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanupCommand)
}
