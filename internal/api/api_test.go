package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/cluster/fake"
	modulefake "github.com/twosson/kubeapt/internal/module/fake"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAPI_routes(t *testing.T) {
	cases := []struct {
		path            string
		expectedCode    int
		expectedContent string
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
			expectedCode: http.StatusNotFound,
		},
		{
			path:            "/content/module/",
			expectedCode:    http.StatusOK,
			expectedContent: "root",
		},
		{
			path:            "/content/module/nested",
			expectedCode:    http.StatusOK,
			expectedContent: "module",
		},
		{
			path:         "/missing",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("GET: %s", tc.path)
		t.Run(name, func(t *testing.T) {
			nsClient := fake.NewNamespaceClient()
			srv := New("/", nsClient)

			m := modulefake.NewModule("module")
			err := srv.RegisterModule(m)
			require.NoError(t, err)

			ts := httptest.NewServer(srv.Handler())
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			u.Path = tc.path

			res, err := http.Get(u.String())
			require.NoError(t, err)
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			if tc.expectedContent != "" {
				assert.Equal(t, tc.expectedContent, string(data))
			}
			assert.Equal(t, tc.expectedCode, res.StatusCode)
		})
	}

}
