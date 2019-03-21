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
	"github.com/homeport/pina-golada/internal/golada/logger"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/homeport/pina-golada/internal/golada/builder"
	"github.com/homeport/pina-golada/pkg/annotation"
	"github.com/homeport/pina-golada/pkg/inspector"

	"github.com/spf13/cobra"
)

var (
	// GoFileSelector is the regex to select go files
	GoFileSelector = regexp.MustCompile("(.*)\\.go")

	parserType string
)

// TestingInterface is just an interfaces
type TestingInterface interface {
	// ProvideSampleGoApp returns the go sample app
	ProvideSampleGoApp() []byte
}

var generateCommand = &cobra.Command{
	Use:   "generate",
	Short: "Generates all requested stubs",
	Long: `Golada will iterate over each file in the current directory and sub-directory. 
	If it finds an interface requesting and defining an asset provider, it will generate one`,
	Run: func(c *cobra.Command, args []string) {
		path := "."
		var (
			parser annotation.Parser
			l      logger.Logger
		)

		switch strings.ToLower(parserType) {
		case "csv":
			parser = annotation.NewCsvParser()
		case "property":
			parser = annotation.NewPropertyParser()
		default:
			log.Fatal("could not find parser for parser type " + parserType)
		}

		if verbose {
			l = logger.NewDefaultLogger(os.Stdout, logger.Debug)
		} else {
			l = logger.NewDefaultLogger(os.Stdout, logger.Info)
		}

		Generate(path, parser, l)
	},
}

// Generate generates the pgl implementations for the given path
func Generate(path string, parser annotation.Parser, logger logger.Logger) {
	logger.Debug("Gray{Debug➤ Generator is now using} LimeGreen{%s} Gray{parser}",
		reflect.TypeOf(parser).Elem().Name())

	fileStream, e := inspector.NewFileStream(path)
	if e != nil {
		log.Fatalf("could not create file stream of directory %s due to %s", path, e.Error())
	}

	astStream := inspector.NewAstStream(fileStream.Filter(func(file inspector.File) bool {
		return GoFileSelector.Match([]byte(file.FileInfo.Name()))
	}))
	for _, i := range astStream.Find() {
		interfaceAnnotation := &builder.PinaGoladaInterface{}
		e := parser.Parse(i.Docs.Text(), interfaceAnnotation)
		if e == annotation.ErrNoAnnotation {
			continue
		} else if e != nil {
			panic(e)
		}

		logger.Info("Aqua{%s}➤ Generating asset provider for LimeGreen{%s}\n", "Pina-Golada",
			i.File.OSFile.FileInfo.Name()+"#"+i.Name.Name)
		logger.Debug("Gray{Debug➤ Found} LimeGreen{%s} Gray{interface}", i.Name.Name)

		output, e := builder.NewBuilder(i, interfaceAnnotation, parser, logger).BuildFile()
		if e != nil {
			log.Fatalf("Could not build file due to " + e.Error())
		}
		baseFileName := filepath.Base(i.Name.Name) + ".go"
		implementationFileName := filepath.Join(filepath.Dir(i.File.OSFile.Path), "pgl"+baseFileName)
		if err := ioutil.WriteFile(implementationFileName, output, os.ModePerm); err != nil {
			log.Fatal("could not create output file due to " + err.Error())
		}
	}
}

func init() {
	generateCommand.PersistentFlags().StringVar(&parserType, "parser", "property", "parser [csv,property]")
	generateCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "--verbose|-v")
	rootCmd.AddCommand(generateCommand)
}
