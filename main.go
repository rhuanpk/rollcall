package main

import (
	"os"
	"os/signal"
	"rollcall/internal/recorder"
	"rollcall/internal/server"
	"syscall"
)

func main() {
	defer server.Listener.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit, os.Interrupt,
		syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGTSTP,
	)

	server.Exec()

	<-quit
	recorder.Process()
}
