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

package cmd

import (
	"bufio"
	"github.com/homeport/pina-golada/internal/golada/logger"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/homeport/pina-golada/internal/golada/builder"
)

var cleanupCommand = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleans all the generated files",
	Long:  "Cleanup will iterate over every go file in the project and delete it if it does start with PinaGolada's identifier string",
	Run: func(c *cobra.Command, args []string) {
		var l logger.Logger
		if verbose {
			l = logger.NewDefaultLogger(os.Stdout, logger.Debug)
		} else {
			l = logger.NewDefaultLogger(os.Stdout, logger.Info)
		}

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
					l.Debug("Gray{Debug➤ Removed pina-golada file} LimeGreen{%s}", file.Name())
				}
			}
			return nil
		})
		if e != nil {
			log.Fatalf("could not iterrate over the files in this directory %s", e.Error())
		}
		l.Info("Aqua{%s}➤ Cleaned project from pina-golada file", "Pina-Golada")
	},
}

func init() {
	cleanupCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "--verbose|-v. Adding this flag to the command will print debug messages for further inside.")
	rootCmd.AddCommand(cleanupCommand)
}
