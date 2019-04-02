// MIT License
//
// Copyright (c) 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package annotation

import (
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"
)

// LabelParser defines a go-style label parser implementation
type LabelParser struct {
}

// NewLabelParser creates a new instance of the label parser
func NewLabelParser() *LabelParser {
	return &LabelParser{}
}

// Parse parses an annotation following the format of the LabelParser
func (c *LabelParser) Parse(comment string, annotation Annotation) (e error) {
	r := regexp.MustCompile(`(.*\n)*( )*\+` + annotation.GetIdentifier() + ` (.*)(\n.*)|($)`)
	result := r.FindAllStringSubmatch(comment, -1)
	if len(result) < 1 {
		return ErrNoAnnotation
	}

	foundAnnotation := make(map[string]string)
	s := result[0][3] // We want the fourth group

	keyValuePair := strings.Split(s, " ")
	for _, paired := range keyValuePair {
		keyToValue := strings.Split(paired, ",")
		if len(keyValuePair) >= 2 {
			foundAnnotation[keyToValue[0]] = keyToValue[1]
		}
	}

	out, e := yaml.Marshal(&foundAnnotation)
	if e != nil {
		return e
	}
	removedQuotes := strings.Replace(string(out), "\"", "", -1)

	if err := yaml.Unmarshal([]byte(removedQuotes), annotation); err != nil {
		return err
	}
	return nil
}
