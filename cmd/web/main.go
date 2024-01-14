package main

import "net/http"

func main() {
	mux := routes()

	_ = http.ListenAndServe(":6969", mux)
}
