package test

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

// TestLoad test the load capacity of server.
//
// Deprecated: need refactor.
func TestLoad(t *testing.T) {
	const clients = 1000
	var wg sync.WaitGroup

	for i := range clients {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			conn, err := net.Dial("tcp4", ":9999")
			if err != nil {
				t.Errorf("client %d: dial error: %v", i, err)
				return
			}
			defer conn.Close()

			conn.Read([]byte{})                 // read title
			conn.Read([]byte{})                 // read setup
			fmt.Fprintln(conn, "1")             // choose option 1
			conn.Read([]byte{})                 // read menu
			fmt.Fprintln(conn, "1")             // choose option 1
			conn.Read([]byte{})                 // read name
			fmt.Fprintf(conn, "client %d\n", i) // send client name
			conn.Read([]byte{})                 // read response

			t.Logf("client %d done\n", i)
		}(i)
	}

	wg.Wait()
}
