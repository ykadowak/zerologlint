package a

import "github.com/rs/zerolog/log"

func f() {
	// The pattern can be written in regular expression.
	var err error
	log.Error().Err(err) // want "missing Msg call for zerolog log method"
}
