package nuhttp

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/nutils"
)

const MimeTypeTextPlain = "text/plain"
const MimeTypeApplicationOctet = "application/octet-stream"

const (
	Http200 = iota + 200
	Http201
)
const (
	Http400 = iota + 400
	Http401
	Http402
	Http403
	Http404
)

type Response struct {
	protocol         string
	httpResponseType int
	header           Header
	body             string
	shouldCompress   bool
}

func responseTypeToString(code int) string {
	switch code {
	case Http200:
		return "200 OK"
	case Http201:
		return "201 Created"
	case Http400:
		return "400 Bad Request"
	case Http404:
		return "404 Not Found"
	}
	return ""
}

func isEncodingTypeValid(encoding string) bool {
	switch strings.Trim(encoding, "\r\n\t ") {
	case "gzip":
		return true
	}
	return false
}

func parseEncodingTypes(value string) string {
	vals := strings.Split(value, ",")
	fmt.Println(vals)
	matchedEncodings := nutils.Filter(vals, isEncodingTypeValid)

	fmt.Println(matchedEncodings)

	// in the future, pick one based off of some configurations/available code
	// in the meantime, we'll just pick gzip
	if len(matchedEncodings) > 0 {
		return strings.Trim(matchedEncodings[0], "\r\n\t ")
	}

	return ""
}

// todo: refactor
func attachOptionalHeaders(headers *[]headerValue, requestHeaders Header) Response {
	response := Response{}
	shouldCompress := false
	for _, header := range requestHeaders.Values {
		if header.name == "Accept-Encoding" {
			encoding := parseEncodingTypes(header.Value)
			if len(encoding) > 0 {
				*headers = append(*headers, headerValue{"Content-Encoding", encoding})
				shouldCompress = true
			}
		}
	}
	response.header.Values = *headers
	response.shouldCompress = shouldCompress
	return response
}

func (r *Response) setContentLength(length int) {
	r.header.SetHeaderValue("Content-Length", fmt.Sprint(length))
}

func (r Response) writeHeaders() string {
	response := r.protocol + " " + responseTypeToString(r.httpResponseType) + "\r\n"
	for _, hdr := range r.header.Values {
		response += (hdr.name + ": " + hdr.Value + "\r\n")
	}
	response += "\r\n"
	return response
}

func (r Response) toString() string {
	r.setContentLength(len(r.body))
	response := r.writeHeaders()
	if len(r.body) > 0 {
		response += r.body
	}
	return response
}

func (r Response) compress() []byte {
	fmt.Println("compressing!")
	if len(r.body) <= 0 {
		r.setContentLength(0)
		response := r.writeHeaders()
		return []byte(response)
	}

	var buf bytes.Buffer

	// gzip needs to be closed before we can read touch the buffer
	compressor := gzip.NewWriter(&buf)
	_, err := compressor.Write([]byte(r.body))
	if err != nil {
		compressor.Close()
		log.Fatal(err)
	}
	compressor.Close()

	bytes := buf.Bytes()
	fmt.Printf("%x\n", bytes)
	r.setContentLength(len(bytes))
	response := r.writeHeaders()
	fmt.Println("size: ", len(bytes))
	return append([]byte(response), bytes...)
}

func (r Response) GetAsBytes() []byte {
	if r.shouldCompress {
		return r.compress()
	}
	return []byte(r.toString())
}

func Ok(protocol string, contentType string, body string, r Request) Response {
	headers := []headerValue{{"Content-Type", contentType}}
	response := attachOptionalHeaders(&headers, r.Header)
	response.protocol = protocol
	response.httpResponseType = Http200
	response.body = body
	return response
}

func Created(protocol string) Response {
	return Response{protocol, Http201, Header{headerPath{}, []headerValue{}}, "", false}
}

func BadRequest(protocol string, err string) Response {
	return Response{protocol, Http400, Header{headerPath{}, []headerValue{}}, err, false}
}

func NotFound(protocol string) Response {
	return Response{protocol, Http404, Header{headerPath{}, []headerValue{}}, "", false}
}
