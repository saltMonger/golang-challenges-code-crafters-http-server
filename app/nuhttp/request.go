package nuhttp

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type headerPath struct {
	verb  string
	Path  string
	proto string
}

type Header struct {
	Path   headerPath
	Values []headerValue
}

func (h Header) HasHeader(input string) bool {
	for _, head := range h.Values {
		if head.name == input {
			return true
		}
	}
	return false
}

func (h Header) GetHeader(input string) (*headerValue, error) {
	for _, header := range h.Values {
		if header.name == input {
			return &header, nil
		}
	}
	return nil, errors.New("no header for: " + input)
}

type Request struct {
	Header Header
	Body   string
}

func (r Request) GetContentSize() (int, error) {
	ret := 0
	for _, header := range r.Header.Values {
		if header.name == "Content-Length" {
			val, err := strconv.Atoi(header.Value)
			if err != nil {
				return 0, err
			}
			ret = val
			break
		}
	}
	return ret, nil
}

func parseHeaders(list []string) []headerValue {
	ret := make([]headerValue, 0)
	for _, input := range list {
		pair := strings.Split(input, ":")
		if len(pair) != 2 {
			continue
		}
		ret = append(ret, headerValue{pair[0], strings.Trim(pair[1], "\r\n\t ")})
	}
	return ret
}

func readBody(input string, contentLen int) string {
	if contentLen <= 0 {
		return ""
	}
	return input[:contentLen]
}

func Parse(input string) Request {
	fmt.Println(input)
	requestLines := strings.Split(input, "\r\n")
	if len(requestLines) < 2 {
		log.Fatal("Malformed request")
	}

	// first line is location/protocol
	// first 3 lines are header
	pathLines := strings.Split(requestLines[0], " ")
	headerPath := headerPath{pathLines[0], strings.Trim(pathLines[1], " \r\n\t"), pathLines[2]}
	headerValues := parseHeaders(requestLines[1 : len(requestLines)-1])

	request := Request{Header{headerPath, headerValues}, ""}
	size, err := request.GetContentSize()
	if err != nil {
		log.Println(err)
	}
	readBody := readBody(requestLines[len(requestLines)-1], size)
	request.Body = readBody
	return request
}
