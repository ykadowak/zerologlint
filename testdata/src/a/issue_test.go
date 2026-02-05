package a

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func MyGoodHandler(w http.ResponseWriter, r *http.Request) {
	evt := log.Ctx(r.Context()).Debug()
	if resource := r.Header.Get("x-resource"); resource != "" {
		evt = evt.Str("resource", resource)
	}

	evt.
		Msg("request received")
}

func MyBadHandler(w http.ResponseWriter, r *http.Request) {
	evt := log.Ctx(r.Context()).Debug()
	if resource := r.Header.Get("x-resource"); resource != "" {
		evt = evt.Str("resource", resource)
	}

	evt.
		Str("method", r.Method).
		Msg("request received")
}
