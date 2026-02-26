package server

import (
	"log"
	"net"
	"rollcall/internal/lists"
	"rollcall/internal/recorder"
	"rollcall/internal/rollcall"
	"strings"
)

// Listener is the global server listener to be closed at end of the main.
var Listener net.Listener

func init() {
	var err error

	recorder.Exec()
	lists.Exec()

	Listener, err = net.Listen("tcp4", ":9999") // ("tcp", "0.0.0.0:9999"), ("tcp4", "")
	if err != nil {
		log.Fatalln("error starting listener:", err)
	}
	log.Println("listening on:", Listener.Addr())
}

// Exec start the TCP server.
func Exec() {
	go func() {
		for {
			conn, err := Listener.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "closed network connection") {
					log.Println("closening the server")
					break
				}

				log.Println("error listener accept:", err)
				continue
			}
			defer conn.Close()
			log.Println(conn.RemoteAddr(), "connected")

			go rollcall.Exec(conn, recorder.File)
		}
	}()
}
