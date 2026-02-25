package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/go-hl/normalize"
)

var mu sync.Mutex

func close(conn net.Conn) {
	log.Println(conn.RemoteAddr(), "disconnected")
	conn.Write([]byte("> disconnected\n"))
	conn.Close()
}

func title() []byte {
	const (
		msg       = "Sistema de Chamada"
		msgLen    = len(msg)
		msgHalf   = msgLen / 2
		msgSub    = msgLen + (msgHalf / 2)
		msgPart   = msgLen + msgHalf
		msgDouble = msgLen * 2
	)

	var title string
	title += strings.Repeat("#", msgPart) + "\n"
	title += fmt.Sprintf("%*s\n", msgSub, msg)
	title += strings.Repeat("#", msgPart) + "\n"

	return []byte(title)
}

func setup(conn net.Conn, reader *bufio.Reader, ln *bool) {
	var setup string
	setup += "1. Terminal\n"
	setup += "2. Aplicativo\n"
	setup += "Opção: "

	conn.Write([]byte(setup))

	input, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			close(conn)
			return
		}
		log.Println(conn.RemoteAddr(), "error reading input:", err)
		return
	}
	input = strings.TrimRight(input, "\r\n")

	option, err := strconv.Atoi(input)
	if err != nil {
		log.Println(conn.RemoteAddr(), "error parsing input:", err)
		return
	}

	if option == 2 {
		// ln = &[]bool{true}[0]
		if ln == nil {
			ln = new(bool)
		}
		*ln = true
	}

	if ln != nil && *ln {
		conn.Write([]byte("\n"))
	}
}

func menu(conn net.Conn, reader *bufio.Reader, ln bool) (int, string) {
	var (
		option int
		input  string
		menu   string
		err    error
	)

	menu += "1. Presença\n"
	menu += "0. Sair\n"
	menu += "Opção: "

	for {
		conn.Write([]byte(menu))

		input, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				close(conn)
				return -1, ""
			}
			log.Println(conn.RemoteAddr(), "error reading input:", err)
			conn.Write([]byte("> error reading input\n"))
			continue
		}
		input = strings.TrimRight(input, "\r\n")

		if ln {
			conn.Write([]byte("\n"))
		}

		option, err = strconv.Atoi(input)
		if err != nil {
			log.Println(conn.RemoteAddr(), "error parsing input:", err)
			conn.Write([]byte("> error parsing input\n"))
			continue
		}

		break
	}

	return option, input
}

func options(conn net.Conn, file *os.File, reader *bufio.Reader, option int, input string, ln bool) {
	mu.Lock()
	defer mu.Unlock()

	switch option {
	case 1:
		conn.Write([]byte("Nome: "))
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println(conn.RemoteAddr(), "error reading input:", err)
			conn.Write([]byte("> error reading input\n"))
			break
		}
		input = strings.TrimRight(input, "\r\n")

		if ln {
			conn.Write([]byte("\n"))
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			log.Println(conn.RemoteAddr(), "error seeking file:", err)
			conn.Write([]byte("> internal server error\n"))
			close(conn)
			break
		}

		bytes, err := io.ReadAll(file)
		if err != nil {
			log.Println(conn.RemoteAddr(), "error reading file:", err)
			conn.Write([]byte("> internal server error\n"))
			close(conn)
			break
		}

		name := normalize.String(input)
		if regexp.MustCompile(fmt.Sprintf(`(?m)^%s$`, name)).Match(bytes) {
			log.Println(conn.RemoteAddr(), "name already present:", name)
			conn.Write([]byte("> name already present\n"))
			close(conn)
			break
		}

		if _, err := file.WriteString(name + "\n"); err != nil {
			log.Println(conn.RemoteAddr(), "error recording rollcall:", err)
			conn.Write([]byte("> error recording rollcall\n"))
			break
		}

		log.Println(conn.RemoteAddr(), "success recording rollcall:", name)
		conn.Write([]byte("> success recording rollcall\n"))

		close(conn)
	case 0:
		close(conn)
	default:
		log.Println(conn.RemoteAddr(), "unknow option:", input)
		conn.Write([]byte("> unknow option\n"))
	}
}

func rollcall(conn net.Conn, file *os.File) {
	defer conn.Close()

	var ln bool
	reader := bufio.NewReader(conn)

	conn.Write(title())
	setup(conn, reader, &ln)

	for {
		if _, err := conn.Read([]byte{}); err != nil {
			if !strings.Contains(err.Error(), "closed network connection") {
				log.Println("error read conn:", err)
			}
			break
		}

		option, input := menu(conn, reader, ln)
		if option == -1 {
			break
		}

		options(conn, file, reader, option, input, ln)
	}
}
