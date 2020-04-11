package overview

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandler_routes(t *testing.T) {
	cases := []struct {
		name         string
		path         string
		values       url.Values
		generator    generator
		expectedCode int
		expectedBody string
	}{
		{
			name:         "GET dynamic content",
			path:         "/api/real",
			values:       url.Values{"namespace": []string{"default"}},
			generator:    newStubbedGenerator(),
			expectedCode: http.StatusOK,
			expectedBody: `{"contents":[{"namespace":"default","type":"real"}]}`,
		},
		{
			name:         "GET stubbed content",
			path:         "/api/stubbed",
			generator:    newStubbedGenerator(),
			expectedCode: http.StatusOK,
			expectedBody: `{"contents":[{"type":"stubbed"}]}`,
		},
		{
			name:         "GET invalid path",
			path:         "/api/missing",
			generator:    newStubbedGenerator(),
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":{"code":404,"message":"content not found"}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newHandler("/api", tc.generator, stubStream)

			ts := httptest.NewServer(h)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			u.Path = tc.path
			u.RawQuery = tc.values.Encode()

			resp, err := http.Get(u.String())
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedCode, resp.StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
			assert.Equal(t, tc.expectedBody, strings.TrimSpace(string(body)))
		})
	}
}

var (
	stubStream = func(ctx context.Context, w http.ResponseWriter, ch chan []byte) {}
)

type stubbedGenerator struct{}

func newStubbedGenerator() *stubbedGenerator {
	return &stubbedGenerator{}
}

func (g *stubbedGenerator) Generate(path, prefix, namespace string) ([]Content, error) {
	switch {
	case strings.HasPrefix(path, "/stubbed"):
		return []Content{
			map[string]string{
				"type": "stubbed",
			},
		}, nil

	case strings.HasPrefix(path, "/real"):
		return []Content{
			map[string]string{
				"type":      "real",
				"namespace": namespace,
			},
		}, nil

	default:
		return nil, contentNotFound
	}
}
