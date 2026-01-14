package validation

import (
	"io"
	"net/http"
	"strings"
)

const (
	MaxBodySize = 5 * 1024 * 1024

	MinBodySize = 100
)

type Validator struct {
	maxBodySize int64
	minBodySize int64
}

func NewValidator() *Validator {
	return &Validator{
		maxBodySize: MaxBodySize,
		minBodySize: MinBodySize,
	}
}

type ValidatedResponse struct {
	Response *http.Response
	Body     []byte
}

func (v *Validator) Validate(resp *http.Response) (*ValidatedResponse, error) {

	if err := v.validateStatus(resp); err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")
	if err := v.validateContentType(contentType); err != nil {
		return nil, err
	}

	if resp.ContentLength > 0 {
		if err := v.validateContentLength(resp.ContentLength); err != nil {
			return nil, err
		}
	}

	body, err := v.readBodyWithLimit(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := v.validateBodySize(int64(len(body))); err != nil {
		return nil, err
	}

	utf8Body, err := DetectAndConvertEncoding(body, contentType)
	if err != nil {
		return nil, err
	}

	return &ValidatedResponse{
		Response: resp,
		Body:     utf8Body,
	}, nil
}

func (v *Validator) validateStatus(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ErrInvalidStatus(resp.StatusCode)
	}
	return nil
}

func (v *Validator) validateContentType(contentType string) error {

	if contentType == "" {
		return nil
	}

	if !strings.Contains(strings.ToLower(contentType), "text/html") {
		return ErrInvalidContentType(contentType)
	}

	return nil
}

func (v *Validator) validateContentLength(length int64) error {
	if length > v.maxBodySize {
		return ErrBodyTooLarge(length)
	}
	return nil
}

func (v *Validator) validateBodySize(size int64) error {
	if size < v.minBodySize {
		return ErrBodyEmpty
	}
	if size > v.maxBodySize {
		return ErrBodyTooLarge(size)
	}
	return nil
}

func (v *Validator) readBodyWithLimit(body io.ReadCloser) ([]byte, error) {

	limitedReader := io.LimitReader(body, v.maxBodySize+1)

	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	if int64(len(data)) > v.maxBodySize {
		return nil, ErrBodyTooLarge(int64(len(data)))
	}

	return data, nil
}
