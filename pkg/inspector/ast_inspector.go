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

// Package inspector hosts the interfaces which analyses existing go source files
// and fetches interface commands and variable setup.
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
	Nodes   [][]ast.Node
	FileMap []*File
}

// Find returns the interfaces declared in the given AstStream
func (a *AstStream) Find() []*AstInterface {
	interfaces := make([]*AstInterface, 0)
	for i, nodeArrayInFile := range a.Nodes { // Loop over all file linked nodes

		genDecl := make([]*ast.GenDecl, 0)

		// This is the file node found in the array. We will only find one because we iterate over nodes of one file
		fileNode := &AstFile{
			OSFile: *a.FileMap[i],
		}

		for _, node := range nodeArrayInFile { // Loop over all nodes in one file
			switch nodeType := node.(type) {
			case *ast.TypeSpec:
				if nodeType.Name.IsExported() {
					switch declaredType := nodeType.Type.(type) {
					case *ast.InterfaceType:
						interfaces = append(interfaces, &AstInterface{
							InterfaceReference: declaredType,
							Name:               nodeType.Name,
							File:               fileNode,
						})
					}
				}
			case *ast.GenDecl:
				genDecl = append(genDecl, nodeType)
			case *ast.File:
				fileNode.FileReference = nodeType
			}
		}

		for _, generalDeclaration := range genDecl {
			for _, interfaceDeclaration := range interfaces {
				if generalDeclaration.End() == interfaceDeclaration.InterfaceReference.End() {
					interfaceDeclaration.Docs = generalDeclaration.Doc
				}
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
	File               *AstFile
}

// AstVariable is a variable that represents a variable
type AstVariable struct {
	ValueReference *ast.ValueSpec
	Docs           *ast.CommentGroup
}

// AstFile represents a file node of the ast package
type AstFile struct {
	FileReference *ast.File
	OSFile        File
}

// NewAstStream creates a new ast stream instance
func NewAstStream(fileStream *FileStream) (astStream *AstStream) {
	nodes := make([][]ast.Node, 0)
	fileSet := token.NewFileSet()

	fileMap := make([]*File, 0)
	fileStream.ForEach(func(file File) {
		parsedFile, _ := parser.ParseFile(fileSet, file.Path, nil, parser.ParseComments)
		fileNodeSet := make([]ast.Node, 0)

		ast.Inspect(parsedFile, func(node ast.Node) bool {
			fileNodeSet = append(fileNodeSet, node)
			return true
		})

		nodes = append(nodes, fileNodeSet)
		fileMap = append(fileMap, &file)
	})

	return &AstStream{
		Nodes:   nodes,
		FileMap: fileMap,
	}
}
