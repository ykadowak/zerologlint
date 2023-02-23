package a

import "github.com/rs/zerolog/log"

func f() {
	// The pattern can be written in regular expression.
	var gopher int // want "identifier is gopher"
	print(gopher)  // want "identifier is gopher"
	var err error
	log.Error().Err(err) // want "missing Msg() or Send()"
}
