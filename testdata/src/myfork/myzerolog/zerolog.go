// Package myzerolog is a minimal zerolog fork stub for testing purposes.
package myzerolog

import "io"

// Event represents a log event (mirrors zerolog.Event).
type Event struct{}

func (e *Event) Str(key, val string) *Event  { return e }
func (e *Event) Int(key string, val int) *Event { return e }
func (e *Event) Err(err error) *Event        { return e }
func (e *Event) Send()                       {}
func (e *Event) Msg(msg string)              {}
func (e *Event) Msgf(format string, v ...interface{}) {}

// Logger mirrors zerolog.Logger.
type Logger struct{}

func New(w io.Writer) Logger { return Logger{} }
func (l Logger) Info() *Event  { return &Event{} }
func (l Logger) Error() *Event { return &Event{} }
func (l Logger) Debug() *Event { return &Event{} }
func (l Logger) Warn() *Event  { return &Event{} }
