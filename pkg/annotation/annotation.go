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

package annotation

import (
	"errors"
)

var (
	// ErrTooManyAnnotations is the error that is thrown when more than one annotation of the same type is declared. The struct will still be filled with the first declared annotation
	ErrTooManyAnnotations = errors.New("the annotation type was declared multiple times in the comment")

	// ErrNoAnnotation is the error that is thrown when no annotations were found in the comment
	ErrNoAnnotation = errors.New("the requested annotation type was not find in the comment")
)

// Parser is an interface that is capable of parsing an annotation based on a string.
type Parser interface {
	// Parse parses the comment into the annotation interface.
	// It will search for the identifier in the comment and parse the entire annotation behind it
	// The method will

	Parse(comment string, annotation Annotation) error
}

// Annotation is a basic interface all annotations have to follow, as they need to provide their identifier
type Annotation interface {
	// GetIdentifier returns the identifier of the annotation
	GetIdentifier() string
}
