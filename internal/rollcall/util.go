package rollcall

import (
	l "log"
	"net"
)

func log(conn net.Conn, values ...any) {
	l.Println(append([]any{conn.RemoteAddr()}, values...)...)
}

func write(conn net.Conn, text string) {
	conn.Write([]byte("> " + text + "\n"))
}

func close(conn net.Conn) {
	log(conn, "disconnected")
	write(conn, "disconnected")

	conn.Close()
}
