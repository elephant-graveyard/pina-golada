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

package logger

import (
	"fmt"
	"github.com/homeport/gonvenience/pkg/v1/bunt"
	"io"
	"strings"
)

// LogLevel defines the type of level the logger will log at
type LogLevel uint32

const (
	// Info defines the info level which will always be printed
	Info LogLevel = iota

	// Debug defines the debug level which will be printed
	Debug
)

// Logger defines a logger that prints strings to the general output
//
// Log logs the message with the given log level
//
// Info logs the message on the info level
//
// Debug logs the message on the debug level
type Logger interface {
	Log(level LogLevel, message string, a ...interface{})
	Info(message string, a ...interface{})
	Debug(message string, a ...interface{})
}

// DefaultLogger defines the default logger implementation
type DefaultLogger struct {
	outputStream    io.Writer
	highestLogLevel LogLevel
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger(outputStream io.Writer, highestLogLevel LogLevel) *DefaultLogger {
	return &DefaultLogger{outputStream: outputStream, highestLogLevel: highestLogLevel}
}

// Log logs the message with the given log level
func (d *DefaultLogger) Log(level LogLevel, message string, a ...interface{}) {
	if d.highestLogLevel >= level {
		message := bunt.Sprintf(message, a...)
		if !strings.HasSuffix(message , "\n") {
			message += "\n"
		}

		fmt.Fprint(d.outputStream, message)
	}
}

// Info logs the message on the info level
func (d *DefaultLogger) Info(message string, a ...interface{}) {
	d.Log(Info, message, a...)
}

// Debug logs the message on the debug level
func (d *DefaultLogger) Debug(message string, a ...interface{}) {
	d.Log(Debug, message, a...)
}
