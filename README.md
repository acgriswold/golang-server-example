# golang-server-example
Simple golang server for creating requests, retrieving/uploading files, etc.


## Get started

```bash
go run ./cmd/golang-server-example # setup and run server
go run ./cmd/golang-server-example --directory files # run server that points to local files (useful for POST files request)
go run ./cmd/golang-server-example --simple # run server with specific handler package (net/http go package on FALSE, hand built request parsing on TRUE) 
```

```bash
curl -v localhost:4221/ # GET the status code of 200
curl -v localhost:4221/echo/<message> # GET the <message> from the request path
curl -v localhost:4221/user-agent # GET the User-Agent from the request header


curl -v localhost:4221/files/<file-name> # GET the contents of <file-name> from the server (need to setup --directory flag and physical directory location)
curl -X localhost:4221/files/<file-name> -H "Content-Type: application/octet-stream" -d "<file-content>" # POST the <file-content> to <file-name> on the server (need to setup --directory flag and physical directory location)
```