package a

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// Test case for issue: False positive with logic
// Both handlers should NOT trigger a warning because they properly dispatch the event

func MyGoodHandler(w http.ResponseWriter, r *http.Request) {
	evt := log.Ctx(r.Context()).Debug()
	if resource := r.Header.Get("x-resource"); resource != "" {
		evt = evt.Str("resource", resource)
	}

	evt.
		Msg("request received")
}

func MyGoodHandlerWithChainedCalls(w http.ResponseWriter, r *http.Request) {
	// This used to trigger a false positive before the fix
	// The event is properly dispatched even though it's chained after conditional reassignment
	evt := log.Ctx(r.Context()).Debug()
	if resource := r.Header.Get("x-resource"); resource != "" {
		evt = evt.Str("resource", resource)
	}

	evt.
		Str("method", r.Method).
		Msg("request received")
}

// Additional test cases to ensure the fix works correctly

func conditionalWithChainedDispatch1() {
	evt := log.Info()
	if true {
		evt = evt.Str("a", "b")
	}
	evt.Str("c", "d").Send()
}

func conditionalWithChainedDispatch2() {
	evt := log.Info()
	if true {
		evt = evt.Str("error", "yes")
	}
	evt.Str("status", "ok").Msg("done")
}

func conditionalWithMultipleChains() {
	evt := log.Debug()
	if true {
		evt = evt.Str("x", "y")
	}
	evt.
		Str("a", "b").
		Str("c", "d").
		Msgf("formatted %s", "message")
}

// This should still trigger a warning - no dispatch
func conditionalWithChainedButNoDispatch() {
	evt := log.Info() // want "must be dispatched by Msg or Send method"
	if true {
		evt = evt.Str("a", "b")
	}
	evt.Str("c", "d")
}
