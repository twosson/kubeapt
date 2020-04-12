package fake

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/apt"
	"github.com/twosson/kubeapt/internal/log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestModule_ContentPath(t *testing.T) {
	m := NewModule("module", log.NopLogger())
	assert.Equal(t, "/module", m.ContentPath())
}

func TestModule_Handler(t *testing.T) {
	m := NewModule("module", log.NopLogger())

	ts := httptest.NewServer(m.Handler("/module"))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	require.NoError(t, err)

	cases := []struct {
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			path:           "/module",
			expectedStatus: http.StatusOK,
			expectedBody:   "root",
		},
		{
			path:           "/module/nested",
			expectedStatus: http.StatusOK,
			expectedBody:   "module",
		},
		{
			path:           "/module/missing",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("GET %s", tc.path)
		t.Run(name, func(t *testing.T) {
			u.Path = tc.path
			resp, err := http.Get(u.String())
			require.NoError(t, err)

			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatus, resp.StatusCode)

			data, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, string(data))
			}
		})
	}
}

func TestModule_Navigation(t *testing.T) {
	m := NewModule("module", log.NopLogger())

	expected := &apt.Navigation{
		Path:  "/module",
		Title: "module",
	}

	got, err := m.Navigation("/module")
	require.NoError(t, err)

	assert.Equal(t, expected, got)
}
