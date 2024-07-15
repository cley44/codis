package main

import (
	"codis"
	"syscall"

	"github.com/samber/do/v2"
)

func main() {
	injector := codis.RegisterAll()

	app := do.MustInvoke[*codis.HTTPAppService](injector)

	go app.ListenAndServe()

	println("Application started")
	_, err := injector.ShutdownOnSignals(syscall.SIGINT, syscall.SIGTERM)
	if err != nil {
		println(err)
	}
	println("Application stopped")
}
