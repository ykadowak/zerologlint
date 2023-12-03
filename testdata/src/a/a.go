package a

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func positives() {
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

	// conditional1
	logger2 := log.Info() // want "must be dispatched by Msg or Send method"
	if err != nil {
		logger2 = log.Error() // want "must be dispatched by Msg or Send method"
	}
	logger2.Str("foo", "bar")

	// conditional2
	loggerCond2 := log.Info().Str("a", "b") // want "must be dispatched by Msg or Send method"
	if err != nil {
		loggerCond2 = loggerCond2.Str("c", "d")
	}
	loggerCond2.Str("foo", "bar")

	// conditional3
	loggerCond3 := log.Info().Str("a", "b") // want "must be dispatched by Msg or Send method"
	if err != nil {
		loggerCond3 = loggerCond3.Str("c", "d")
	}
	if err != nil {
		loggerCond3 = loggerCond3.Str("e", "f")
	}
	loggerCond3.Str("foo", "bar")

	// conditional4
	var event *zerolog.Event
	if true {
		event = log.Info() // want "must be dispatched by Msg or Send method"
	} else {
		event = log.Warn() // want "must be dispatched by Msg or Send method"
	}
	if true {
		event = event.Err(nil)
	}
	event.Str("foo", "bar")

	// defer patterns
	defer log.Info() // want "must be dispatched by Msg or Send method"

	logger3 := log.Error() // want "must be dispatched by Msg or Send method"
	defer logger3.Err(err).Str("foo", "bar").Int("foo", 1)

	defer log.Info(). // want "must be dispatched by Msg or Send method"
				Str("foo", "bar").
				Dict("dict", zerolog.Dict().
					Str("bar", "baz").
					Int("n", 1),
		)

	// logger instance
	logger4 := zerolog.New(os.Stdout)
	logger4.Info() // want "must be dispatched by Msg or Send method"

	// custom object marshaller
	f := &Foo{Bar: &Bar{}}
	log.Info().Object("foo", f) // want "must be dispatched by Msg or Send method"

	// logger instance not dispatched within other function
	l := log.Info() // want "must be dispatched by Msg or Send method"
	badDispatcher(l)
}

func negatives() {
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

	// conditional2
	loggerCond2 := log.Info().Str("a", "b")
	if err != nil {
		loggerCond2 = loggerCond2.Str("c", "d")
	}
	loggerCond2.Str("foo", "bar")
	loggerCond2.Send()

	// conditional3
	loggerCond3 := log.Info().Str("a", "b")
	if err != nil {
		loggerCond3 = loggerCond3.Str("c", "d")
	}
	if err != nil {
		loggerCond3 = loggerCond3.Str("e", "f")
	}
	loggerCond3.Str("foo", "bar")
	loggerCond3.Send()

	// conditional4
	var event *zerolog.Event
	if true {
		event = log.Info()
	} else {
		event = log.Warn()
	}
	if true {
		event = event.Err(nil)
	}
	event.Send()

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

	// logger instance
	logger4 := zerolog.New(os.Stdout)
	logger4.Info().Send()

	// custom object marshaller
	f := &Foo{Bar: &Bar{}}
	log.Info().Object("foo", f).Msg("")

	// zerolog.Event dispatched within other function
	l := log.Info()
	goodDispatcher(l)
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

func badDispatcher(e *zerolog.Event) {
	e.Str("foo", "bar")
}

func goodDispatcher(e *zerolog.Event) {
	e.Send()
}
