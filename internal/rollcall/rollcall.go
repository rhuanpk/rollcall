package rollcall

import (
	"bufio"
	"net"
	"os"
)

// Exec execute the rollcall system.
func Exec(conn net.Conn, file *os.File) {
	defer conn.Close()

	var line bool
	reader := bufio.NewReader(conn)

	conn.Write(title())
	setup(conn, reader, &line)
	rollcall(conn, reader, file, line)
}
