package nuhttp

import (
	"fmt"
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
	headers          []headerValue
	body             string
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

func attachOptionalHeaders(headers *[]headerValue, requestHeaders Header) *[]headerValue {
	for _, header := range requestHeaders.Values {
		if header.name == "Accept-Encoding" {
			encoding := parseEncodingTypes(header.Value)
			if len(encoding) > 0 {
				*headers = append(*headers, headerValue{"Content-Encoding", encoding})
			}
		}
	}
	return headers
}

func (r Response) ToString() string {
	response := r.protocol + " " + responseTypeToString(r.httpResponseType) + "\r\n"
	for _, hdr := range r.headers {
		response += (hdr.name + ": " + hdr.Value + "\r\n")
	}
	response += "\r\n"
	if len(r.body) > 0 {
		response += r.body
	}
	return response
}

func Ok(protocol string, contentType string, body string, r Request) Response {
	headers := []headerValue{{"Content-Type", contentType}, {"Content-Length", fmt.Sprint(len(body))}}
	headers = *attachOptionalHeaders(&headers, r.Header)
	return Response{protocol, Http200, headers, body}
}

func Created(protocol string) Response {
	return Response{protocol, Http201, []headerValue{}, ""}
}

func BadRequest(protocol string, err string) Response {
	return Response{protocol, Http400, make([]headerValue, 0), err}
}

func NotFound(protocol string) Response {
	return Response{protocol, Http404, make([]headerValue, 0), ""}
}
