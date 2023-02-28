package a

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func bad() {
	// The pattern can be written in regular expression.
	log.Fatal() // want "missing Msg or Send call for zerolog log method"
	log.Panic() // want "missing Msg or Send call for zerolog log method"
	log.Debug() // want "missing Msg or Send call for zerolog log method"
	log.Info() // want "missing Msg or Send call for zerolog log method"
	log.Warn() // want "missing Msg or Send call for zerolog log method"
	log.Error() // want "missing Msg or Send call for zerolog log method"

	var err error
	log.Error().Err(err) // want "missing Msg or Send call for zerolog log method"
	log.Error().Err(err).Str("foo", "bar").Int("foo", 1) // want "missing Msg or Send call for zerolog log method"
	log.Info(). // want "missing Msg or Send call for zerolog log method"
    Str("foo", "bar").
    Dict("dict", zerolog.Dict().
        Str("bar", "baz").
        Int("n", 1),
    )
}

func ok() {
	log.Fatal().Send()
	log.Panic().Msg("")
	log.Debug().Send()
	log.Info().Msg("")
	log.Warn().Send()
	log.Error().Msg("")

	var err error
	log.Error().Err(err).Send()
	log.Error().Err(err).Str("foo", "bar").Int("foo", 1).Msg("")
	log.Info().
    Str("foo", "bar").
    Dict("dict", zerolog.Dict().
        Str("bar", "baz").
        Int("n", 1),
    ).Send()

	// FIXME: this should pass. Use SSA
	// logger := log.Error()
	// logger.Send()
}
