package validation

import (
	"bytes"
	"io"
	"unicode/utf8"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/unicode"
)

func DetectAndConvertEncoding(body []byte, contentType string) ([]byte, error) {

	e, name, _ := charset.DetermineEncoding(body, contentType)

	if e == unicode.UTF8 || name == "utf-8" {
		return body, nil
	}

	decoder := e.NewDecoder()
	utf8Body, err := io.ReadAll(decoder.Reader(bytes.NewReader(body)))
	if err != nil {
		return nil, ErrInvalidEncoding(name)
	}

	return utf8Body, nil
}

func IsValidUTF8(body []byte) bool {
	return utf8.Valid(body)
}
