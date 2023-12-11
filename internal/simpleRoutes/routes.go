package simpleRoutes

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"sync"

	"os"
	"strings"

	"github.com/acgriswold/golang-server-example/internal/handler"
)

func HandleConnection(conn net.Conn, wg *sync.WaitGroup, fileDirectory string) {
	defer wg.Done()
	defer conn.Close()

	request := extractRequest(conn)

	switch request.Method {
	case "GET":
		handleGet(conn, &request, fileDirectory)
	case "POST":
		handlePost(conn, &request, fileDirectory)
	default:
		handleResponse(buildResponse("405 Method Not Allowed"), conn)
	}
}

func handleGet(conn net.Conn, request *Req, fileDirectory string) {
	if request.Path == "/" {
		handleResponse(buildResponse("200 OK"), conn)
	} else if strings.HasPrefix(request.Path, "/echo") {
		body := strings.TrimPrefix(strings.TrimPrefix(request.Path, "/echo"), "/")
		handleResponse(buildResponseWithBody("200 OK", "text/plain", body), conn)
	} else if strings.HasPrefix(request.Path, "/user-agent") {
		handleResponse(buildResponseWithBody("200 OK", "text/plain", request.Headers["User-Agent"]), conn)
	} else if strings.HasPrefix(request.Path, "/files/") {
		handleFileResponse(request.Path, conn, fileDirectory)
	} else {
		handleResponse(buildResponse("404 Not Found"), conn)
	}
}

func handlePost(conn net.Conn, request *Req, fileDirectory string) {
	if strings.HasPrefix(request.Path, "/files/") {
		fileName := strings.Split(request.Path, "/files/")[1]
		path := filepath.Join(fileDirectory, fileName)

		_, err := os.Stat(path)
		handler.CheckFileError(err, "Error checking status of file")

		writeError := os.WriteFile(path, []byte(request.Body), 0644)
		handler.Check(writeError, "Error writing file to disk", false)

		handleFileResponse(request.Path, conn, fileDirectory)
	} else {
		handleResponse(buildResponse("404 Not Found"), conn)
	}
}

func buildResponse(status string) []byte {
	return []byte("HTTP/1.1 " + status + "\r\n\r\n")
}

func buildResponseWithBody(status string, contentType string, body string) []byte {
	response := "HTTP/1.1 " + status + "\r\n" +
		"Content-Type: " + contentType + "\r\n" +
		"Content-Length: " + strconv.Itoa(len(body)) +
		"\r\n\r\n" + body

	fmt.Println("response ", response)

	return []byte(response)
}

func handleResponse(response []byte, conn net.Conn) {
	_, err := conn.Write(response)
	handler.Check(err, "Error writing data on connection", true)
}

func handleFileResponse(filePath string, conn net.Conn, fileDirectory string) {
	fileName := strings.Split(filePath, "/files/")[1]
	path := filepath.Join(fileDirectory, fileName)
	file, err := os.ReadFile(path)

	if err != nil {
		handleResponse(buildResponse("404 Not Found"), conn)
	} else {
		handleResponse(buildResponseWithBody("200 OK", "application/octet-stream", string(file)), conn)
	}
}

func extractRequest(conn net.Conn) Req {
	buffer := make([]byte, 4096)

	byteSize, err := conn.Read(buffer)
	handler.Check(err, "Failed to read contents of HTTP Request", true)

	fmt.Print("request ", string(buffer[:byteSize]))

	content := strings.Split(string(buffer[:byteSize]), "\r\n")
	request := strings.Split(content[0], " ")

	body := strings.TrimSpace(content[len(content)-1])

	headers := make(map[string]string)
	for _, line := range content[1:] {
		if strings.TrimSpace(line) == "" {
			break
		}

		kv := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		headers[key] = value
	}

	return Req{
		Method:  request[0],
		Path:    request[1],
		Version: request[2],

		Headers: headers,
		Body:    body,
	}
}

type Req struct {
	Method  string
	Path    string
	Version string
	Body    string
	Headers map[string]string
}
