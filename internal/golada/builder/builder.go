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
	"github.com/homeport/pina-golada/internal/golada/logger"
	"go/format"
	"path/filepath"
	"strings"

	"github.com/homeport/pina-golada/pkg/annotation"
	"github.com/homeport/pina-golada/pkg/compressor"
	"github.com/homeport/pina-golada/pkg/files"
	"github.com/homeport/pina-golada/pkg/generator"
	"github.com/homeport/pina-golada/pkg/inspector"
)

var (
	// IdentifierString is pasted in the first line of each file and will mark a
	// file as a pina golada generated file. It adheres to the Go rules to mark
	// files as generated:
	// https://github.com/golang/go/issues/13560#issuecomment-288457920
	IdentifierString = "// Code generated by pina-golada (https://github.com/homeport/pina-golada). DO NOT EDIT."

	// InternalDecompressMethod defines the method that is generated by pina-golada to decompress
	// a given hex string. It cannot be defined by an interface
	InternalDecompressMethod = "requestAssetByPath"
)

// PinaGoladaInterface is the struct used for the pina golada interface annotation
type PinaGoladaInterface struct {
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
	logger              logger.Logger
}

// NewBuilder Creates a new builder instance
func NewBuilder(target *inspector.AstInterface,
	interfaceAnnotation *PinaGoladaInterface,
	parser annotation.Parser,
	l logger.Logger) *Builder {
	return &Builder{target: target, interfaceAnnotation: interfaceAnnotation, parser: parser, logger: l}
}

// BuildFile Builds the file
func (b Builder) BuildFile() (by []byte, err error) {
	if len(b.interfaceAnnotation.Injector) < 1 {
		return nil, errors.New("no injector variable defined on " + b.target.Name.Name)
	}

	goGenerator := generator.NewFileGoGenerator(b.target.File.FileReference.Name.Name)
	structName := "PGL" + b.target.Name.Name
	receiverVariableName := "p"
	receiverType := receiverVariableName + " *" + structName

	goGenerator.Struct(structName, func(s generator.StructGenerator) {})
	goGenerator.Import("bytes")
	goGenerator.Import("encoding/hex")
	goGenerator.Import("fmt")
	goGenerator.Import("github.com/homeport/pina-golada/pkg/compressor")
	goGenerator.Import("github.com/homeport/pina-golada/pkg/files")

	goGenerator.Method(InternalDecompressMethod, func(method generator.MethodGenerator) {
		// As we generate a method that uses fmt.Errorf therefore it has to ignore that we don't pass a value
		// to the Sprintf method, therefore we pass nil
		method.Receiver(receiverType).Parameters("compressorType string", "hexedContent string").
			ReturnTypes("files.Directory", "error").
			Body([]string{
				fmt.Sprintf("c := compressor.DefaultRegistry.Find(compressorType)"),
				fmt.Sprintf(`if c == nil{return nil, fmt.Errorf("could not find compressor for %s", compressorType)}`, "%s"),
				fmt.Sprintf(`decodedBytes , er := hex.DecodeString(hexedContent)`),
				fmt.Sprintf(`if er != nil {return nil , er}`),
				fmt.Sprintf(`return c.Decompress(bytes.NewBuffer(decodedBytes))`),
			}...)
	})

	for _, method := range b.target.InterfaceReference.Methods.List {
		if len(method.Names) < 1 {
			return nil, errors.New("method commented with " + method.Doc.Text() + " has no name")
		}
		methodName := method.Names[0].Name
		if methodName == "requestAssetByPath" {
			return nil, fmt.Errorf("interface %s defined a method with the name %s which is internally used",
				b.target.Name.Name,
				InternalDecompressMethod)
		}

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

		directory := files.NewRootDirectory()
		isDir, e := files.LoadFromDiskAndType(directory, methodAnnotation.Asset)
		if e != nil {
			return nil, e
		}

		files.WalkFileTree(directory, func(file files.File) {
			b.logger.Debug("Gray{Debug➤ Found asset file} White{%s}", file.Name().String())
		})

		compressorType := compressor.DefaultRegistry.Find(methodAnnotation.Compressor)
		if compressorType == nil {
			return nil, errors.New("could not find compressor for " + methodAnnotation.Compressor)
		}
		buffer := &bytes.Buffer{}
		if err := compressorType.Compress(directory, buffer); err != nil {
			return nil, err
		}

		b.logger.Debug("Gray{Debug➤ Compressed asset} LimeGreen{%s} Gray{for method} LimeGreen{%s} "+
			"Gray{to} LimeGreen{%d} Gray{bytes}",
			methodAnnotation.Asset, b.target.Name.Name+"#"+methodName, buffer.Len())

		goGenerator.Method(methodName, func(method generator.MethodGenerator) {
			assetProviderCall := fmt.Sprintf("%s.%s(\"%s\", \"%s\")", receiverVariableName,
				InternalDecompressMethod,
				methodAnnotation.Compressor,
				hex.EncodeToString(buffer.Bytes()))

			method.Receiver(receiverType).ReturnTypes("files.Directory", "error")

			if isDir {
				goGenerator.Import("github.com/homeport/pina-golada/pkg/files/paths")
				method.Body(fmt.Sprintf("dir, decompressError := %s", assetProviderCall))
				method.Body(fmt.Sprintf(`if decompressError != nil {return dir, decompressError}`))
				method.Body(fmt.Sprintf(`return dir.Directory(paths.Of("%s")) , nil`,
					filepath.Base(methodAnnotation.Asset)))
			} else {
				method.Body(fmt.Sprintf("return %s", assetProviderCall))
			}
		})
	}

	goGenerator.Method("init", func(method generator.MethodGenerator) {
		method.Body(b.interfaceAnnotation.Injector + " = &" + structName + "{}")
	})
	b.logger.Debug("Gray{Debug➤ Generated initialize method for} LimeGreen{%s}", b.target.Name.Name)

	outputBuffer := &bytes.Buffer{}
	outputBuffer.WriteString(IdentifierString + "\n\n")
	goGenerator.Flush(outputBuffer)
	return format.Source(outputBuffer.Bytes())
}
