package nuhttp

import "fmt"

const Http200 = 200
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
	case Http400:
		return "400 Bad Request"
	case Http404:
		return "404 Not Found"
	}
	return ""
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

func Ok(protocol string, body string) Response {
	fmt.Println("body: " + body)
	headers := []headerValue{{"Content-Type", "text/plain"}, {"Content-Length", fmt.Sprint(len(body))}}
	return Response{protocol, Http200, headers, body}
}

func BadRequest(protocol string, err string) Response {
	return Response{protocol, Http400, make([]headerValue, 0), err}
}

func NotFound(protocol string) Response {
	return Response{protocol, Http404, make([]headerValue, 0), ""}
}
