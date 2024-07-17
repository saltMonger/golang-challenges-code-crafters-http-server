package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/file"
	"github.com/codecrafters-io/http-server-starter-go/app/nuhttp"
)

const Http200 = 200
const (
	Http400 = iota + 400
	Http401
	Http402
	Http403
	Http404
)

const readChunkSize = 1024

var directoryVal string
var fileDir file.FileDirectory

func initialize() {
	flag.StringVar(&directoryVal, "directory", "/tmp/", "file serving directory")
	flag.Parse()
}

func parseChunk(c net.Conn) (int, []byte, error) {
	var received int
	// TODO: actually implement chunk sizing
	buffer := bytes.NewBuffer(nil)
	for {
		chunk := make([]byte, readChunkSize)
		read, err := c.Read(chunk)
		if err != nil {
			return received, buffer.Bytes(), err
		}
		received += read
		buffer.Write(chunk[:read])

		if read == 0 || read < readChunkSize {
			break
		}
	}
	return received, buffer.Bytes(), nil
}

func routeRequest(r nuhttp.Request) nuhttp.Response {
	path := strings.Split(r.Header.Path.Path, "/")
	if len(path) == 2 && len(path[1]) == 0 {
		return nuhttp.Ok("HTTP/1.1", nuhttp.MimeTypeTextPlain, "", r)
	}

	if path[1] == "echo" {
		return nuhttp.Ok("HTTP/1.1", nuhttp.MimeTypeTextPlain, path[2], r)
	}

	if path[1] == "user-agent" {
		body, err := r.Header.GetHeader("User-Agent")
		if err != nil {
			return nuhttp.BadRequest("HTTP/1.1", err.Error())
		}
		return nuhttp.Ok("HTTP/1.1", nuhttp.MimeTypeTextPlain, body.Value, r)
	}

	if path[1] == "files" {

		if r.Header.Path.Verb == "POST" {
			err := fileDir.CreateFile(path[2], r.Body)
			if err != nil {
				return nuhttp.BadRequest("HTTP/1.1", err.Error())
			}
			return nuhttp.Created("HTTP/1.1")

		} else if r.Header.Path.Verb == "GET" {
			file, err := fileDir.GetFile(path[2])
			if err != nil {
				return nuhttp.NotFound("HTTP/1.1")
			}
			return nuhttp.Ok("HTTP/1.1", nuhttp.MimeTypeApplicationOctet, string(file), r)
		}
	}

	return nuhttp.NotFound("HTTP/1.1")
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	_, data, err := parseChunk(conn)
	if err != nil {
		log.Fatal(err)
	}

	requestString := string(data)
	request := nuhttp.Parse(requestString)
	response := routeRequest(request)
	written, err := io.WriteString(conn, response.ToString())
	if err != nil {
		log.Println(err)
	}
	fmt.Println(response.ToString())
	fmt.Printf("Bytes written: %d\n", written)
}

func main() {
	initialize()
	fileDir = file.MakeDirectory(directoryVal)

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go func(c net.Conn) {
			handleClient(c)
		}(conn)
	}
}
