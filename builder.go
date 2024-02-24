package reqbuilder

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// BodyMarshaler is an interface for marshaling request bodies.
type BodyMarshaler interface {
	Marshal(v any) ([]byte, error)
}

var _ BodyMarshaler = BodyMarshalerFunc(nil)

// BodyMarshalerFunc is a function type that implements the BodyMarshaler interface.
type BodyMarshalerFunc func(v any) ([]byte, error)

// BodyMarshalerFunc is a function type that implements the BodyMarshaler interface.
func (fn BodyMarshalerFunc) Marshal(v any) ([]byte, error) {
	return fn(v)
}

// URLValuesMarshaler is an interface for marshaling URL parameters.
type URLValuesMarshaler interface {
	Marshal(v any) (url.Values, error)
}

var _ URLValuesMarshaler = URLValuesMarshalerFunc(nil)

// URLValuesMarshalerFunc is a function type that implements the URLValuesMarshaler interface.
type URLValuesMarshalerFunc func(v any) (url.Values, error)

// Marshal marshals the given interface into URL values.
func (fn URLValuesMarshalerFunc) Marshal(v any) (url.Values, error) {
	return fn(v)
}

// Builder is a struct for configuring and building HTTP requests.
type Builder struct {
	method             string
	path               string
	baseURL            string
	body               any
	header             http.Header
	bodyMarshaler      BodyMarshaler
	urlValuesMarshaler URLValuesMarshaler
	params             any
}

// NewBuilder creates a new instance of Builder with the provided base URL and options.
func NewBuilder(baseURL string, opts ...BuilderOption) Builder {
	options := defaultBuilderOptions(opts...)

	return Builder{
		baseURL:            strings.TrimSuffix(baseURL, " /"),
		method:             http.MethodGet,
		header:             http.Header{},
		body:               http.NoBody,
		bodyMarshaler:      options.bodyMarshaler,
		urlValuesMarshaler: options.urlValuesMarshaler,
	}
}

// WithMethod sets the HTTP method for the request.
func (b Builder) WithMethod(method string) Builder {
	b.method = method

	return b
}

// WithPath sets the URL path for the request.
func (b Builder) WithPath(path string) Builder {
	b.path = path

	return b
}

// WithBody sets the request body.
func (b Builder) WithBody(body any) Builder {
	b.body = body

	return b
}

// WithBodyMarshaler sets the marshaler for the request body.
func (b Builder) WithBodyMarshaler(bodyMarshaler BodyMarshaler) Builder {
	b.bodyMarshaler = bodyMarshaler

	return b
}

// WithHeaders sets the request headers.
func (b Builder) WithHeaders(headers ...http.Header) Builder {
	mergeHeaders(b.header, headers...)

	return b
}

// WithParams sets the URL parameters.
func (b Builder) WithParams(params any) Builder {
	b.params = params

	return b
}

// Build constructs the HTTP request based on the configured parameters.
func (b Builder) Build(ctx context.Context) (*http.Request, error) {
	body, err := b.getBodyReader()
	if err != nil {
		return nil, fmt.Errorf("getBodyReader: %w", err)
	}

	urlValues, err := b.urlValuesMarshaler.Marshal(b.params)
	if err != nil {
		return nil, fmt.Errorf("marshal url values: %w", err)
	}

	rawURL, err := buildURL(b.baseURL, b.path, urlValues.Encode())
	if err != nil {
		return nil, fmt.Errorf("buildURL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, b.method, rawURL, body)
	if err != nil {
		return nil, fmt.Errorf("newRequestWithContext: %w", err)
	}

	replaceHeaders(req.Header, b.header)

	return req, nil
}

func (b Builder) getBodyReader() (io.Reader, error) {
	switch body := b.body.(type) {
	case nil:
		return http.NoBody, nil

	case []byte:
		return bytes.NewReader(body), nil

	case string:
		return strings.NewReader(body), nil

	case io.Reader:
		return body, nil
	}

	data, err := b.bodyMarshaler.Marshal(b.body)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	return bytes.NewBuffer(data), nil
}
