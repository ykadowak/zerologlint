package a

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func bad() {
	log.Error() // want "missing to dispatch with Msg or Send function. nothing will be logged"
	log.Info()  // want "missing to dispatch with Msg or Send function. nothing will be logged"
	log.Fatal() // want "missing to dispatch with Msg or Send function. nothing will be logged"
	log.Debug() // want "missing to dispatch with Msg or Send function. nothing will be logged"
	log.Warn()  // want "missing to dispatch with Msg or Send function. nothing will be logged"

	var err error
	log.Error().Err(err)                                 // want "missing to dispatch with Msg or Send function. nothing will be logged"
	log.Error().Err(err).Str("foo", "bar").Int("foo", 1) // want "missing to dispatch with Msg or Send function. nothing will be logged"

	logger := log.Error() // want "missing to dispatch with Msg or Send function. nothing will be logged"
	logger.Err(err).Str("foo", "bar").Int("foo", 1)

	// include zerolog.Dict()
	log.Info(). // want "missing to dispatch with Msg or Send function. nothing will be logged"
			Str("foo", "bar").
			Dict("dict", zerolog.Dict().
				Str("bar", "baz").
				Int("n", 1),
		)

	// conditional
	logger2 := log.Info() // want "missing to dispatch with Msg or Send function. nothing will be logged"
	if err != nil {
		logger2 = log.Error() // want "missing to dispatch with Msg or Send function. nothing will be logged"
	}
	logger2.Str("foo", "bar")
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
