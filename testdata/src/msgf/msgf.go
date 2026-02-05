package msgf

import (
	"github.com/rs/zerolog/log"
)

func testMsgf() {
	// When checkmsgf flag is enabled, these should be flagged as errors
	log.Info().Msgf("test") // want "must be dispatched by Msg or Send method"

	log.Error().Str("key", "value").Msgf("error: %s", "msg") // want "must be dispatched by Msg or Send method"

	// These should always be acceptable
	log.Info().Msg("test")
	log.Info().Send()
	log.Info().MsgFunc(func() string { return "test" })
}
