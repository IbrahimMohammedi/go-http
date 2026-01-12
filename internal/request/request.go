package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	fullRequest := string(data)

	if len(fullRequest) == 0 {
		return nil, fmt.Errorf("empty request")
	}

	lines := strings.Split(fullRequest, "\r\n")

	reqLine, err := ParseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *reqLine,
	}, nil
}

func ParseRequestLine(line string) (*RequestLine, error) {
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line: expected 3 parts, got %d", len(parts))
	}

	method := parts[0]
	target := parts[1]
	versionRaw := parts[2]

	switch method {
	case "GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS":
	default:
		return nil, fmt.Errorf("invalid method: %s", method)
	}

	if !strings.HasPrefix(target, "/") {
		return nil, fmt.Errorf("invalid request target: %s", target)
	}

	if !strings.HasPrefix(versionRaw, "HTTP/") {
		return nil, fmt.Errorf("invalid version format: %s", versionRaw)
	}

	versionNumber := strings.TrimPrefix(versionRaw, "HTTP/")

	if versionNumber != "1.1" {
		return nil, fmt.Errorf("invalid version number: %s", versionNumber)
	}

	return &RequestLine{
		HttpVersion:   versionNumber,
		RequestTarget: target,
		Method:        method,
	}, nil
}
