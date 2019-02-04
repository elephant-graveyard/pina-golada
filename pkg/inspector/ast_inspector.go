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

package inspector

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// AstFilter defines an ast filter type
type AstFilter func(n *ast.Node) bool

// AstStream is a stream of nodes that can be inspected and transformed
type AstStream struct {
	Nodes   []ast.Node
	FileMap []*File
}

// Find returns the interfaces declared in the given AstStream
func (a *AstStream) Find() []*AstInterface {
	interfaces := make([]*AstInterface, 0)
	genDecl := make([]*ast.GenDecl, 0)

	for i, node := range a.Nodes {
		switch nodeType := node.(type) {
		case *ast.TypeSpec:
			if nodeType.Name.IsExported() {
				switch declaredType := nodeType.Type.(type) {
				case *ast.InterfaceType:
					interfaces = append(interfaces, &AstInterface{
						InterfaceReference: declaredType,
						Name:               nodeType.Name,
						File:               *a.FileMap[i],
					})
				}
			}
		case *ast.GenDecl:
			genDecl = append(genDecl, nodeType)
		}
	}

	for _, generalDeclaration := range genDecl {
		for _, interfaceDeclaration := range interfaces {
			if generalDeclaration.End() == interfaceDeclaration.InterfaceReference.End() {
				interfaceDeclaration.Docs = generalDeclaration.Doc
			}
		}
	}

	return interfaces
}

// AstInterface is a simple struct describing an interface
type AstInterface struct {
	InterfaceReference *ast.InterfaceType
	Docs               *ast.CommentGroup
	Name               *ast.Ident
	File               File
}

// AstVariable is a variable that represents a variable
type AstVariable struct {
	ValueReference *ast.ValueSpec
	Docs           *ast.CommentGroup
}

// NewAstStream creates a new ast stream instance
func NewAstStream(fileStream *FileStream) (astStream *AstStream) {
	nodes := make([]ast.Node, 0)
	fileSet := token.NewFileSet()

	fileMap := make([]*File, 0)
	fileStream.ForEach(func(file File) {
		parsedFile, _ := parser.ParseFile(fileSet, file.Path, nil, parser.ParseComments)
		ast.Inspect(parsedFile, func(node ast.Node) bool {
			nodes = append(nodes, node)
			fileMap = append(fileMap, &file)
			return true
		})
	})

	return &AstStream{
		Nodes:   nodes,
		FileMap: fileMap,
	}
}
