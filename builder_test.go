package reqbuilder_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/go-reqbuilder"
)

func TestBuilder_BuildWithHeader(t *testing.T) {
	t.Parallel()

	builder := reqbuilder.NewBuilder("https://example.com")

	t.Run("headers are not persisted across multiple builds", func(t *testing.T) {
		t.Parallel()

		req, err := builder.
			WithHeaders(reqbuilder.JSONContentHeader).
			Build(context.Background())
		require.NoError(t, err)
		require.Equal(t, req.Header, http.Header{"Content-Type": []string{"application/json"}})

		req, err = builder.Build(context.Background())
		require.NoError(t, err)
		require.Empty(t, req.Header)
	})
}

func TestBuilder_Build(t *testing.T) {
	t.Parallel()

	type Body struct {
		Field1 string `json:"field1,omitempty"`
		Field2 string `json:"field2,omitempty"`
		Field3 string `json:"field3,omitempty"`
	}

	type Params struct {
		Param1 string   `url:"param1,omitempty"`
		Param2 int      `url:"param2,omitempty"`
		Param3 []string `url:"param3,omitempty"`
	}

	tests := []struct {
		name    string
		builder reqbuilder.Builder
		want    *http.Request
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "with body",
			builder: reqbuilder.
				NewBuilder("https://example.com").
				WithMethod(http.MethodPost).
				WithPath("/foo/bar").
				WithHeaders(reqbuilder.JSONAcceptHeader, reqbuilder.JSONContentHeader).
				WithBody(Body{
					Field1: "field1",
					Field2: "field2",
					Field3: "field3",
				}),

			want: &http.Request{
				URL:    parseURL(t, "https://example.com/foo/bar"),
				Method: http.MethodPost,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
					"Accept":       []string{"application/json"},
				},
				Body: NewBody(t, Body{
					Field1: "field1",
					Field2: "field2",
					Field3: "field3",
				}),
			},
		},
		{
			name: "with params",
			builder: reqbuilder.
				NewBuilder("https://example.com").
				WithMethod(http.MethodGet).
				WithPath("/foo/bar").
				WithHeaders(reqbuilder.JSONAcceptHeader, reqbuilder.URLEncodedContentHeader).
				WithParams(Params{
					Param1: "param1",
					Param2: 100,
					Param3: []string{"value1", "value2"},
				}),

			want: &http.Request{
				URL:    parseURL(t, "https://example.com/foo/bar?param1=param1&param2=100&param3=value1&param3=value2"),
				Method: http.MethodGet,
				Header: http.Header{
					"Content-Type": []string{"application/x-www-form-urlencoded"},
					"Accept":       []string{"application/json"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.builder.Build(context.Background())

			if tt.wantErr != nil {
				tt.wantErr(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.Method, got.Method)
			require.Equal(t, tt.want.URL, got.URL)
			require.Equal(t, tt.want.Header, got.Header)
			require.Equal(t, readAll(t, tt.want.Body), readAll(t, got.Body))
		})
	}
}

func parseURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()

	u, err := url.Parse(rawURL)
	require.NoError(t, err)

	return u
}

func NewBody(t *testing.T, v any) io.ReadCloser {
	t.Helper()

	switch a := v.(type) {
	case string:
		return io.NopCloser(strings.NewReader(a))
	case []byte:
		return io.NopCloser(bytes.NewReader(a))
	default:
		b, err := json.Marshal(v)
		require.NoError(t, err)

		return io.NopCloser(bytes.NewReader(b))
	}
}

func readAll(t *testing.T, r io.Reader) string {
	t.Helper()

	if r == nil {
		return ""
	}

	b, err := io.ReadAll(r)
	require.NoError(t, err)

	return string(b)
}
