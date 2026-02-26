package main

import (
	"os"
	"os/signal"
	"rollcall/internal/recorder"
	"rollcall/internal/server"
)

func main() {
	defer server.Listener.Close()
	defer recorder.File.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	server.Init()

	<-quit
	recorder.Clean()
}
