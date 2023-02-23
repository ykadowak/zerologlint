package a

import log2 "github.com/rs/zerolog/log"

func importDecl() {
	// The pattern can be written in regular expression.
	var err error
	log2.Error().Err(err) // want "missing Msg call for zerolog log method"
}
