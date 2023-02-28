package a

import (
	"github.com/rs/zerolog"
	log2 "github.com/rs/zerolog/log"
)

func bad_alias() {
	// The pattern can be written in regular expression.
	log2.Fatal() // want "missing Msg or Send call for zerolog log method"
	log2.Panic() // want "missing Msg or Send call for zerolog log method"
	log2.Debug() // want "missing Msg or Send call for zerolog log method"
	log2.Info() // want "missing Msg or Send call for zerolog log method"
	log2.Warn() // want "missing Msg or Send call for zerolog log method"
	log2.Error() // want "missing Msg or Send call for zerolog log method"

	var err error
	log2.Error().Err(err) // want "missing Msg or Send call for zerolog log method"
	log2.Error().Err(err).Str("foo", "bar").Int("foo", 1) // want "missing Msg or Send call for zerolog log method"
	log2.Info(). // want "missing Msg or Send call for zerolog log method"
    Str("foo", "bar").
    Dict("dict", zerolog.Dict().
        Str("bar", "baz").
        Int("n", 1),
    )
}

func ok_alias() {
	log2.Fatal().Send()
	log2.Panic().Msg("")
	log2.Debug().Send()
	log2.Info().Msg("")
	log2.Warn().Send()
	log2.Error().Msg("")

	var err error
	log2.Error().Err(err).Send()
	log2.Error().Err(err).Str("foo", "bar").Int("foo", 1).Msg("")
	log2.Info().
    Str("foo", "bar").
    Dict("dict", zerolog.Dict().
        Str("bar", "baz").
        Int("n", 1),
    ).Send()

	// FIXME: this should pass. Use SSA
	// logger := log2.Error()
	// logger.Send()
}
