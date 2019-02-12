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

package builder

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/homeport/pina-golada/pkg/annotation"
	"github.com/homeport/pina-golada/pkg/compressor"
	"github.com/homeport/pina-golada/pkg/files"
	"github.com/homeport/pina-golada/pkg/generator"
	"github.com/homeport/pina-golada/pkg/inspector"
)

var (
	// IdentifierString is pasted in the first line of each file and will mark a file as a pina golada generated file
	IdentifierString = "//THIS FILE WAS GENERATED USING PINA-GOLADA, PLEASE DO NOT EDIT"
)

// PinaGoladaInterface is the struct used for the pina golada interface annotation
type PinaGoladaInterface struct {
	Package  string `yaml:"package"`
	Injector string `yaml:"injector"`
}

// GetIdentifier returns the identifier of the interface
func (PinaGoladaInterface) GetIdentifier() string {
	return "@pgl"
}

// PinaGoladaMethod is the struct used for the pina golada interface annotation
type PinaGoladaMethod struct {
	Asset        string `yaml:"asset"`
	Compressor   string `yaml:"compressor"`
	AbsolutePath bool   `yaml:"absolute"`
}

// GetIdentifier returns the identifier of the interface
func (PinaGoladaMethod) GetIdentifier() string {
	return "@pgl"
}

// Builder is able to build a file
type Builder struct {
	target              *inspector.AstInterface
	interfaceAnnotation *PinaGoladaInterface
	parser              annotation.Parser
}

// NewBuilder Creates a new builder instance
func NewBuilder(target *inspector.AstInterface, interfaceAnnotation *PinaGoladaInterface, parser annotation.Parser) *Builder {
	return &Builder{target: target, interfaceAnnotation: interfaceAnnotation, parser: parser}
}

// BuildFile Builds the file
func (b Builder) BuildFile() (by []byte, err error) {
	if len(b.interfaceAnnotation.Package) < 1 {
		return nil, errors.New("no package defined on " + b.target.Name.Name)
	}
	if len(b.interfaceAnnotation.Injector) < 1 {
		return nil, errors.New("no injector variable defined on " + b.target.Name.Name)
	}

	goGenerator := generator.NewFileGoGenerator(b.interfaceAnnotation.Package)
	structName := "PGL" + b.target.Name.Name
	receiverType := "p *" + structName

	goGenerator.Struct(structName, func(s generator.StructGenerator) {})
	goGenerator.Import("bytes")
	goGenerator.Import("encoding/hex")
	goGenerator.Import("errors")
	goGenerator.Import("github.com/homeport/pina-golada/pkg/compressor")
	goGenerator.Import("github.com/homeport/pina-golada/pkg/files")

	for _, method := range b.target.InterfaceReference.Methods.List {
		if len(method.Names) < 1 {
			return nil, errors.New("method commented with " + method.Doc.Text() + " has no name")
		}
		methodName := method.Names[0].Name

		methodAnnotation := &PinaGoladaMethod{}
		if err := b.parser.Parse(method.Doc.Text(), methodAnnotation); err != nil {
			return nil, err
		}
		if len(methodAnnotation.Asset) < 1 {
			return nil, errors.New("no asset path was provided for " + methodName)
		}
		if len(methodAnnotation.Compressor) < 1 {
			return nil, errors.New("no compressor was provided for " + methodName)
		}
		if !methodAnnotation.AbsolutePath && !strings.HasPrefix(methodAnnotation.Asset, "."+string(filepath.Separator)) {
			methodAnnotation.Asset = filepath.Join(".", methodAnnotation.Asset)
		}

		isDir, directory, e := loadFromDisk(files.NewRootDirectory(), methodAnnotation.Asset)
		if e != nil {
			return nil, e
		}

		compressorType := compressor.DefaultRegistry.Find(methodAnnotation.Compressor)
		if compressorType == nil {
			return nil, errors.New("could not find compressor for " + methodAnnotation.Compressor)
		}
		buffer := &bytes.Buffer{}
		if err := compressorType.Compress(directory, buffer); err != nil {
			return nil, err
		}

		goGenerator.Method(methodName, func(method generator.MethodGenerator) {
			method.Receiver(receiverType).ReturnTypes("files.Directory", "error").Body([]string{
				fmt.Sprintf("c := compressor.DefaultRegistry.Find(\"%s\")", methodAnnotation.Compressor),
				fmt.Sprintf(`if c == nil{return nil,errors.New("could not find compressor for %s")}`, methodAnnotation.Compressor),
				fmt.Sprintf(`decodedBytes , er := hex.DecodeString("%s")`, hex.EncodeToString(buffer.Bytes())),
				fmt.Sprintf(`if er != nil {return nil , er}`),
			}...)

			if isDir {
				goGenerator.Import("github.com/homeport/pina-golada/pkg/files/paths")
				method.Body(fmt.Sprintf(`dir, decompressError := c.Decompress(bytes.NewBuffer(decodedBytes))`))
				method.Body(fmt.Sprintf(`if decompressError != nil {return dir, decompressError}`))
				method.Body(fmt.Sprintf(`return dir.Directory(paths.Of("%s")) , nil`,
					filepath.Base(methodAnnotation.Asset)))
			} else {
				method.Body(fmt.Sprintf(`return c.Decompress(bytes.NewBuffer(decodedBytes))`))
			}
		})
	}

	goGenerator.Method("init", func(method generator.MethodGenerator) {
		method.Body(b.interfaceAnnotation.Injector + " = &" + structName + "{}")
	})

	outputBuffer := &bytes.Buffer{}
	outputBuffer.WriteString(IdentifierString + "\n")
	goGenerator.Flush(outputBuffer)
	return outputBuffer.Bytes(), nil
}
