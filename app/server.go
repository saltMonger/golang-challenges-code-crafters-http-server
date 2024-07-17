package main

import (
	"bufio"
	"fmt"
	"io"
	"log"

	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)

const Http200 = 200
const (
	Http400 = iota + 400
	Http401
	Http402
	Http403
	Http404
)

type headerPath struct {
	verb  string
	path  string
	proto string
}

type header struct {
	path headerPath
}

type request struct {
	header header
}

func parseRequest(input string) request {
	fmt.Println(input)
	requestLines := strings.Split(input, "\r\n")
	if len(requestLines) < 2 {
		log.Fatal("Malformed request")
	}

	// first 3 lines are header
	pathLines := strings.Split(requestLines[0], " ")
	headerPath := headerPath{pathLines[0], strings.Trim(pathLines[1], " \r\n\t"), pathLines[2]}
	return request{header{headerPath}}
}

func routeRequest(r request) int {
	if r.header.path.path == "/" {
		return Http200
	}
	return Http404
}

func respond(httpResponse int) string {
	response := "HTTP/1.1 "
	switch httpResponse {
	case Http200:
		response += "200 OK"
	case Http404:
		response += "400 Not Found"
	}
	response += "\r\n\r\n"
	return response
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	//var buf bytes.Buffer
	reader := bufio.NewReader(conn)
	var requestString string
	for n := 0; n < 3; n++ {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		requestString += str
	}

	//io.Copy(&buf, conn)
	request := parseRequest(requestString)
	//request := parseRequest(buf.String())
	//fmt.Println("req: ", request.header.path.path)
	//res := bytes.NewBufferString(respond(routeRequest(request)))
	io.WriteString(conn, respond(routeRequest(request)))
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	handleClient(conn)
}
