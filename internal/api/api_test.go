package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/cluster/fake"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAPI_routes(t *testing.T) {
	cases := []struct {
		path         string
		expectedCode int
	}{
		{
			path:         "/namespaces",
			expectedCode: http.StatusOK,
		},
		{
			path:         "/navigation",
			expectedCode: http.StatusOK,
		},
		{
			path:         "/content/",
			expectedCode: http.StatusOK,
		},
		{
			path:         "/content/nested",
			expectedCode: http.StatusOK,
		},
		{
			path:         "/missing",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("GET: %s", tc.path)
		t.Run(name, func(t *testing.T) {
			o := fake.NewSimpleClusterOverview()
			srv := New("/", o)

			ts := httptest.NewServer(srv)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			u.Path = tc.path

			res, err := http.Get(u.String())
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)
		})
	}

}
