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
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// CsvParser is a csv implementation of the annotation parser.
// It will parse the annotation following a csv format, where the key,value pairs are separated by a ,
type CsvParser struct {
}

// NewCsvParser creates a new parser with the default typ registry
func NewCsvParser() *CsvParser {
	return &CsvParser{}
}

// Parse parses an annotation following the format of the CsvParser
func (c *CsvParser) Parse(comment string, annotation Annotation) (e error) {
	r := regexp.MustCompile(`.*@` + annotation.GetIdentifier() + `\((.*)\).*`)
	result := r.FindAllStringSubmatch(comment, -1)
	if len(result) < 1 {
		return ErrNoAnnotation
	}

	foundAnnotation := make(map[string]string)
	s := result[0][1] // We want the second group

	keyValuePair := strings.Split(s, ";")
	for _, paired := range keyValuePair {
		keyToValue := strings.Split(paired, ",")
		if len(keyToValue) >= 2 {
			foundAnnotation[keyToValue[0]] = keyToValue[1]
		}
	}

	out, e := yaml.Marshal(&foundAnnotation)
	if e != nil {
		return e
	}

	removedQuotes := strings.Replace(string(out), "\"", "", -1)
	return yaml.Unmarshal([]byte(removedQuotes), annotation)
}
