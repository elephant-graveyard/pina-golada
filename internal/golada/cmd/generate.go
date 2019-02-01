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

// go:generate echo Hello there

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

var (

	// GoFileSelector is the regex to select go files
	GoFileSelector = regexp.MustCompile("(.*)\\.go")
)

var generateCommand = &cobra.Command{
	Use:   "generate",
	Short: "Generates all requested stubs",
	Long: `Golada will iterate over each file in the current directory and sub-directory. 
	If it finds an interface requesting and defining an asset provider, it will generate one`,
	Run: func(c *cobra.Command, args []string) {
		filepath.Walk(".", func(root string, file os.FileInfo, e error) error {
			if GoFileSelector.Match([]byte(file.Name())) {
				fmt.Println(file.Name())
			}
			return nil
		})
		fmt.Println("a")
	},
}

func init() {
	rootCmd.AddCommand(generateCommand)
}
