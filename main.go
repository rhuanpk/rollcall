package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

const (
	fileName   = "rollcall.txt"
	folderName = "assets"
)

var (
	file     *os.File
	listener net.Listener
)

func init() {
	var err error

	if err := os.Mkdir(folderName, 0775); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatalln("error creating folder:", err)
	}

	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("error opening file:", err)
	}

	if _, err := file.WriteString(time.Now().Format(time.DateTime) + "\n"); err != nil {
		log.Fatalln("error recording datetime:", err)
	}

	listener, err = net.Listen("tcp4", ":9999") // ("tcp", "0.0.0.0:9999"), ("tcp4", "")
	if err != nil {
		log.Fatalln("error starting listener:", err)
	}
	log.Println("listening on:", listener.Addr())
}

func main() {
	defer file.Close()
	defer listener.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "closed network connection") {
					log.Println("closening the server")
					break
				}
				log.Println("error in listener accept:", err)
				continue
			}
			log.Println(conn.RemoteAddr(), "connected")

			go rollcall(conn, file)
		}
	}()

	<-quit
	if err := os.Rename(fileName, filepath.Join(
		folderName, time.Now().Format(
			fmt.Sprintf("%s_%s.txt", time.DateOnly, time.TimeOnly),
		),
	)); err != nil {
		log.Println("error changing rollcall name:", err)
	}
}
