package rollcall

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"rollcall/pkg/errs"
	"strconv"
	"strings"
)

func scan(conn net.Conn, reader *bufio.Reader, safe, line bool) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			close(conn)
			return "", errs.ErrClosedConn
		}

		log(conn, "error reading input:", err)
		write(conn, "error reading input")

		if strings.Contains(err.Error(), errs.ErrInvalidSyntax.Error()) {
			err = errs.ErrInvalidSyntax
		}

		return "", err
	}

	if safe && line {
		conn.Write([]byte("\n"))
	}

	return strings.TrimRight(input, "\r\n"), nil
}

func menu(conn net.Conn, reader *bufio.Reader, exit, safe, line bool, options ...string) int {
	const optionExit = 0
	var option int

	if len(options) <= 0 {
		log(conn, "error missing menu options")
		write(conn, "internal server error")

		close(conn)
		return -1
	}

	var label strings.Builder
	for index, option := range options {
		fmt.Fprintf(&label, "%d. %s\n", index+1, option)
	}
	if exit {
		label.WriteString("0. Sair\n")
	}
	label.WriteString("Opção: ")

	for {
		conn.Write([]byte(label.String()))

		input, err := scan(conn, reader, safe, line)
		if err != nil {
			if errors.Is(err, errs.ErrClosedConn) {
				return -1
			}

			if !safe && errors.Is(err, errs.ErrInvalidSyntax) {
				return -2
			}

			continue
		}

		option, err = strconv.Atoi(input)
		if safe && err != nil {
			log(conn, "error parsing input:", err)
			write(conn, "error parsing input")
			continue
		}

		if exit && option == optionExit {
			close(conn)
			return -1
		}

		if safe && (option < 1 || option > len(options)) {
			log(conn, "unknow option:", input)
			write(conn, "unknow option")
			continue
		}

		return option
	}
}
