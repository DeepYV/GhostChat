package main

import (
	"github.com/bmizerany/pat"
	"net/http"
	"webrtc/internal/handlers"
)

func routes() http.Handler {

	go handlers.ListToWsChannel()
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndPoint))

	return mux

}
