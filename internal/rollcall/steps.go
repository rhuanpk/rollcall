package rollcall

import (
	"bufio"
	"errors"
	"net"
	"regexp"
	"rollcall/internal/lists"
	"rollcall/pkg/errs"
	"strings"
	"sync"
	"time"

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

func rollcall(conn net.Conn, reader *bufio.Reader, line bool) {
	const exit, safe = true, true

	mu.Lock()
	defer mu.Unlock()

	for {
		conn.SetReadDeadline(time.Now().Add(time.Second))
		if _, err := conn.Read([]byte{}); err != nil {
			conn.SetReadDeadline(time.Time{})
			if !strings.Contains(err.Error(), "closed network connection") &&
				!strings.Contains(err.Error(), "i/o timeout") {
				log(conn, "error reading connection:", err)
			}
			return
		}
		conn.SetReadDeadline(time.Time{})

		option := menu(conn, reader, exit, safe, line, "Presença")
		if option == -1 {
			return
		}

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
			name := normalize.String(input)

			var pattern strings.Builder
			pattern.WriteString(`.*`)
			for part := range strings.SplitSeq(name, " ") {
				pattern.WriteString(part + ".*")
			}
			regex := regexp.MustCompile(pattern.String())

			matches := regex.FindAllString(lists.String, -1)
			if len(matches) <= 0 {
				log(conn, "name not found")
				write(conn, "name not found")
				continue
			}
			if len(matches) > 1 {
				log(conn, "many same names")
				write(conn, "many same names")
				continue
			}
			name = regex.FindString(lists.String)

			present, ok := lists.List[name]
			if !ok {
				log(conn, "full name not found")
				write(conn, "internal server error")

				close(conn)
				return
			}

			if present {
				log(conn, "name already present:", name)
				write(conn, "name already present")

				close(conn)
				return
			}

			lists.List[name] = true
			log(conn, "success recording presence:", name)
			write(conn, "success recording presence")

			close(conn)
			return
		}
	}
}
