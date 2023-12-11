package main

import (
	"flag"
	"fmt"
	"net/http"

	"net"
	"sync"

	"github.com/acgriswold/golang-server-example/internal/handler"
	"github.com/acgriswold/golang-server-example/internal/routes"
	"github.com/acgriswold/golang-server-example/internal/simpleRoutes"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	var directoryFlag = flag.String("directory", ".", "directory to serve files from")
	var simpleServerFlag = flag.Bool("simple", false, "determines whether to run simple server or net/http server")
	flag.Parse()

	fmt.Println("flag server simple", *simpleServerFlag)

	if !*simpleServerFlag {
		router := routes.NewRouter()
		err := http.ListenAndServe(":4221", router)
		handler.Check(err, "Failed to bind to port 4221", false)
	} else {
		l, err := net.Listen("tcp", "0.0.0.0:4221")
		handler.Check(err, "Failed to bind to port 4221", false)

		var wg sync.WaitGroup
		for {
			conn, err := l.Accept()
			handler.Check(err, "Error accepting connection", true)

			wg.Add(1)
			go simpleRoutes.HandleConnection(conn, &wg, *directoryFlag)
		}

		wg.Wait()
	}
}
