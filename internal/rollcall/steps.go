package rollcall

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"rollcall/pkg/errs"
	"strings"
	"sync"

	"github.com/go-hl/normalize"
)

var mu sync.Mutex

func setup(conn net.Conn, reader *bufio.Reader, line *bool) {
	const exit, safe = false, false

	option := menu(conn, reader, exit, safe, false, "Terminal", "Aplicativo")
	if option == -1 {
		return
	}

	if option == 2 {
		// ln = &[]bool{true}[0]
		if line == nil {
			line = new(bool)
		}
		*line = true
	}

	if line != nil && *line {
		conn.Write([]byte("\n"))
	}
}

func rollcall(conn net.Conn, reader *bufio.Reader, file *os.File, line bool) {
	const exit, safe = true, true

	for {
		if _, err := conn.Read([]byte{}); err != nil {
			if !strings.Contains(err.Error(), "closed network connection") {
				log(conn, "error reading connection:", err)
			}
			return
		}

		option := menu(conn, reader, exit, safe, line, "Presença")
		if option == -1 {
			return
		}

		mu.Lock()
		defer mu.Unlock()

		switch option {
		case 1:
			conn.Write([]byte("Nome: "))

			input, err := scan(conn, reader, safe, line)
			if err != nil {
				if errors.Is(err, errs.ErrClosedConn) {
					return
				}
				continue
			}

			if _, err := file.Seek(0, io.SeekStart); err != nil {
				log(conn, "error seeking record file:", err)
				write(conn, "internal server error")

				close(conn)
				return
			}

			bytes, err := io.ReadAll(file)
			if err != nil {
				log(conn, "error reading record file:", err)
				write(conn, "internal server error")

				close(conn)
				return
			}

			name := normalize.String(input)
			if regexp.MustCompile(fmt.Sprintf(`(?m)^%s$`, name)).Match(bytes) {
				log(conn, "name already present:", name)
				write(conn, "name already present")

				close(conn)
				return
			}

			if _, err := file.WriteString(name + "\n"); err != nil {
				log(conn, "error recording presence:", err)
				write(conn, "error recording presence")

				close(conn)
				return
			}

			log(conn, "success recording presence:", name)
			write(conn, "success recording presence")

			close(conn)
			return
		}
	}
}
