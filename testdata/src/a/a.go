package a

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func bad() {
	log.Error() // want "missing Msg or Send call for zerolog log method"
	log.Info()  // want "missing Msg or Send call for zerolog log method"
	log.Fatal() // want "missing Msg or Send call for zerolog log method"
	log.Debug() // want "missing Msg or Send call for zerolog log method"
	log.Warn()  // want "missing Msg or Send call for zerolog log method"

	var err error
	log.Error().Err(err)                                 // want "missing Msg or Send call for zerolog log method"
	log.Error().Err(err).Str("foo", "bar").Int("foo", 1) // want "missing Msg or Send call for zerolog log method"

	logger := log.Error() // want "missing Msg or Send call for zerolog log method"
	logger.Err(err).Str("foo", "bar").Int("foo", 1)
}

func ok() {
	log.Fatal().Send()
	log.Panic().Msg("")
	log.Debug().Send()
	log.Info().Msg("")
	log.Warn().Send()
	log.Error().Msg("")

	log.Error().Str("foo", "bar").Send()
	var err error
	log.Error().Err(err).Str("foo", "bar").Int("foo", 1).Msg("")

	logger := log.Error()
	logger.Send()

	// include zerolog.Dict()
	log.Info().
		Str("foo", "bar").
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).Send()

	// conditional
	logger2 := log.Info()
	if err != nil {
		logger2 = log.Error()
	}
	logger2.Send()
}
