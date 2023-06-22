package a

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func bad() {
	log.Error() // want "must be dispatched by Msg or Send method"
	log.Info()  // want "must be dispatched by Msg or Send method"
	log.Fatal() // want "must be dispatched by Msg or Send method"
	log.Debug() // want "must be dispatched by Msg or Send method"
	log.Warn()  // want "must be dispatched by Msg or Send method"

	var err error
	log.Error().Err(err)                                 // want "must be dispatched by Msg or Send method"
	log.Error().Err(err).Str("foo", "bar").Int("foo", 1) // want "must be dispatched by Msg or Send method"

	logger := log.Error() // want "must be dispatched by Msg or Send method"
	logger.Err(err).Str("foo", "bar").Int("foo", 1)

	// include zerolog.Dict()
	log.Info(). // want "must be dispatched by Msg or Send method"
			Str("foo", "bar").
			Dict("dict", zerolog.Dict().
				Str("bar", "baz").
				Int("n", 1),
		)

	// conditional
	logger2 := log.Info() // want "must be dispatched by Msg or Send method"
	if err != nil {
		logger2 = log.Error() // want "must be dispatched by Msg or Send method"
	}
	logger2.Str("foo", "bar")

	// defer patterns
	defer log.Info()      // want "must be dispatched by Msg or Send method"

	logger3 := log.Error() // want "must be dispatched by Msg or Send method"
	defer logger3.Err(err).Str("foo", "bar").Int("foo", 1)

	defer log.Info(). // want "must be dispatched by Msg or Send method"
				Str("foo", "bar").
				Dict("dict", zerolog.Dict().
					Str("bar", "baz").
					Int("n", 1),
		)

	// custom object marshaller
	f := &Foo{Bar: &Bar{}}

	log.Info().Object("foo", f) // want "must be dispatched by Msg or Send method"
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

	// dispatch variation
	log.Info().Msgf("")
	log.Info().MsgFunc(func() string { return "foo" })

	// defer patterns
	defer log.Info().Msg("")

	logger3 := log.Info()
	defer logger3.Msg("")

	defer log.Info().
		Str("foo", "bar").
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).Send()

	// custom object marshaller
	f := &Foo{Bar: &Bar{}}

	log.Info().Object("foo", f).Msg("")
}

type Marshaller interface {
	MarshalZerologObject(event *zerolog.Event)
}

type Foo struct {
	Bar Marshaller
}

func (f *Foo) MarshalZerologObject(event *zerolog.Event) {
	f.Bar.MarshalZerologObject(event)
}

type Bar struct{}

func (b *Bar) MarshalZerologObject(event *zerolog.Event) {
	event.Str("key", "value")
}
