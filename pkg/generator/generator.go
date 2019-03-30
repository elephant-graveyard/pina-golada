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

// Package generator hosts the generator interfaces, which are the base
// for building go source files from runtime.
// The represent go source elements like methods and structs during the runtime.
package generator

import (
	"io"
	"strings"
)

var (
	// EndOfLine is simply the end of line character
	EndOfLine = "\n"
	// Intend contains the tab char
	Intend = "\t"
	// Quote contains a single double quote
	Quote = string('"')
)

// GoGenerator is an interface that is capable of generating go files
//
// Package sets the package name of the generated file
//
// Import adds an import to the file
//
// Method adds a new method with a method body to the go generator
//
// Struct adds a new struct to the generator
//
// Flush flushes the generated content into the writer
type GoGenerator interface {
	Import(importName string) GoGenerator
	Method(name string, consumer func(m MethodGenerator)) GoGenerator
	Struct(name string, consumer func(s StructGenerator)) GoGenerator
	Flush(writer io.Writer)
}

// FileGoGenerator is a basic implementation of the GoGenerator that is capable of creating a go file
type FileGoGenerator struct {
	packageName string
	imports     []string
	methods     []string
	structs     []string
}

// Method adds a new method to the generator
func (f *FileGoGenerator) Method(name string, consumer func(method MethodGenerator)) GoGenerator {
	methodGenerator := &SimpleMethodGenerator{name: name}
	consumer(methodGenerator)

	writer := &strings.Builder{}
	methodGenerator.Flush(writer)

	f.methods = append(f.methods, writer.String())
	return f
}

// NewFileGoGenerator returns a new instance of the FileGoGenerator
func NewFileGoGenerator(packageName string) *FileGoGenerator {
	return &FileGoGenerator{packageName: packageName}
}

// Flush flushes the content to a writer
func (f *FileGoGenerator) Flush(writer io.Writer) {
	builder := strings.Builder{}
	builder.WriteString("package " + f.packageName + EndOfLine)
	builder.WriteString(EndOfLine)
	if len(f.imports) > 0 {
		builder.WriteString("import (" + EndOfLine)
		for _, importString := range f.imports {
			builder.WriteString(Intend + Quote + importString + Quote + EndOfLine)
		}
		builder.WriteString(")" + EndOfLine)
		builder.WriteString(EndOfLine)
	}

	for _, s := range f.structs {
		builder.WriteString(s)
		builder.WriteString(EndOfLine)
	}
	for _, method := range f.methods {
		builder.WriteString(method)
		builder.WriteString(EndOfLine)
	}
	_, _ = writer.Write([]byte(builder.String()))
}

// Import adds a new import to the generator
func (f *FileGoGenerator) Import(importName string) GoGenerator {
	contains := func(list []string, value string) bool {
		for _, entry := range list {
			if entry == value {
				return true
			}
		}

		return false
	}

	if !contains(f.imports, importName) {
		f.imports = append(f.imports, importName)
	}

	return f
}

// Struct adds a new struct to the generator
func (f *FileGoGenerator) Struct(name string, consumer func(s StructGenerator)) GoGenerator {
	structGenerator := &SimpleStructGenerator{name: name}
	consumer(structGenerator)

	stringWriter := &strings.Builder{}
	structGenerator.Flush(stringWriter)

	f.structs = append(f.structs, stringWriter.String())
	return f
}

// StructGenerator is a generator for struct inside a go file
//
// Flush flushes the generator to the writer
//
// Field adds a new field to the struct
type StructGenerator interface {
	Field(name string, fieldType string) StructGenerator
}

// SimpleStructGenerator is a simple implementation of the StructGenerator interface
type SimpleStructGenerator struct {
	name   string
	fields []string
}

// Flush flushes the content of the struct
func (s *SimpleStructGenerator) Flush(writer io.Writer) {
	builder := strings.Builder{}
	builder.WriteString("type " + s.name + " struct {" + EndOfLine)
	for _, field := range s.fields {
		builder.WriteString(Intend + field + EndOfLine)
	}
	builder.WriteString("}" + EndOfLine)
	_, _ = writer.Write([]byte(builder.String()))
}

// Field adds a new field to the struct generator
func (s *SimpleStructGenerator) Field(name string, fieldType string) StructGenerator {
	s.fields = append(s.fields, name+" "+fieldType)
	return s
}

// MethodGenerator defines an object that can generate a go method
//
// Parameters sets the parameters
//
// Receiver adds a receiver string
//
// ReturnTypes adds the return types
//
// Body adds the method body
type MethodGenerator interface {
	Parameters(parameters ...string) MethodGenerator
	Receiver(r string) MethodGenerator
	ReturnTypes(returnTypes ...string) MethodGenerator
	Body(body ...string) MethodGenerator
}

// SimpleMethodGenerator is a basic implementation of the MethodGenerator interface
type SimpleMethodGenerator struct {
	name        string
	parameters  []string
	receiver    string
	returnTypes []string
	body        []string
}

// Parameters sets the parameters of the generated method
func (m *SimpleMethodGenerator) Parameters(parameters ...string) MethodGenerator {
	m.parameters = append(m.parameters, parameters...)
	return m
}

// Receiver sets the receiver type of the generated method
func (m *SimpleMethodGenerator) Receiver(r string) MethodGenerator {
	m.receiver = r
	return m
}

// ReturnTypes sets the return-types of the generated method
func (m *SimpleMethodGenerator) ReturnTypes(returnTypes ...string) MethodGenerator {
	m.returnTypes = append(m.returnTypes, returnTypes...)
	return m
}

// Body sets the body of the generated method
func (m *SimpleMethodGenerator) Body(body ...string) MethodGenerator {
	m.body = append(m.body, body...)
	return m
}

// Flush flushes the content of the generator to a writer instance
func (m *SimpleMethodGenerator) Flush(writer io.Writer) {
	builder := &strings.Builder{}

	// function definition with optional receiver
	builder.WriteString("func ")
	if len(m.receiver) > 0 {
		builder.WriteString("(" + m.receiver + ") ")
	}

	// function name and function parameters
	builder.WriteString(m.name + "(" + strings.Join(m.parameters, ", ") + ")")

	// function return types (if any)
	if len(m.returnTypes) > 0 {
		builder.WriteString(" (" + strings.Join(m.returnTypes, ", ") + ")")
	}

	// function body
	builder.WriteString(" {" + EndOfLine)
	for _, body := range m.body {
		builder.WriteString(Intend + body + EndOfLine)
	}
	builder.WriteString("}" + EndOfLine)

	writer.Write([]byte(builder.String()))
}
