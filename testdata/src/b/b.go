package b

import (
	"myfork/myzerolog"
	"myfork/myzerolog/log"
)

func positives() {
	log.Error() // want "must be dispatched by Msg or Send method"
	log.Info()  // want "must be dispatched by Msg or Send method"

	var err error
	log.Error().Err(err).Str("foo", "bar") // want "must be dispatched by Msg or Send method"

	logger := myzerolog.New(nil)
	logger.Info() // want "must be dispatched by Msg or Send method"
}

func negatives() {
	log.Fatal().Send()
	log.Info().Msg("")
	log.Error().Str("foo", "bar").Send()

	var err error
	log.Error().Err(err).Str("foo", "bar").Msg("")

	logger := myzerolog.New(nil)
	logger.Info().Send()

	_ = err
}
