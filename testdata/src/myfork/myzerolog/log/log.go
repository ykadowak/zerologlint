// Package log provides global logging functions using the myzerolog fork.
package log

import "myfork/myzerolog"

func Info() *myzerolog.Event  { return &myzerolog.Event{} }
func Error() *myzerolog.Event { return &myzerolog.Event{} }
func Debug() *myzerolog.Event { return &myzerolog.Event{} }
func Warn() *myzerolog.Event  { return &myzerolog.Event{} }
func Fatal() *myzerolog.Event { return &myzerolog.Event{} }
