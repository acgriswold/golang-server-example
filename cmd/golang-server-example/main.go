package main

import (
	"flag"
	"fmt"

	"net"
	"sync"

	"github.com/acgriswold/golang-server-example/internal/handler"
	"github.com/acgriswold/golang-server-example/internal/simpleRoutes"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	handler.Check(err, "Failed to bind to port 4221", false)

	var directoryFlag = flag.String("directory", ".", "directory to serve files from")
	flag.Parse()

	var wg sync.WaitGroup
	for {
		conn, err := l.Accept()
		handler.Check(err, "Error accepting connection", true)

		wg.Add(1)
		go simpleRoutes.HandleConnection(conn, &wg, *directoryFlag)
	}

	wg.Wait()
}
