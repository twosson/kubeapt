package dash

import (
	"context"
	"fmt"
	"github.com/twosson/kubeapt/internal/api"
	"github.com/twosson/kubeapt/internal/cluster/fake"
	"github.com/twosson/kubeapt/internal/module"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	modulefake "github.com/twosson/kubeapt/internal/module/fake"
)

func TestDash_Run(t *testing.T) {
	cases := []struct {
		name         string
		hasCustomURL bool
		expected     string
	}{
		{
			name:     "embedded dashboard ui",
			expected: "embedded",
		},
		{
			name:         "custom dashboard ui",
			hasCustomURL: true,
			expected:     "custom",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			namespace := "default"
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)

			var uiURL string
			if tc.hasCustomURL {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "custom")
				}))
				defer ts.Close()
				uiURL = ts.URL
			}

			defaultHandler := func() (http.Handler, error) {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "embedded")
				}), nil
			}

			nsClient := fake.NewNamespaceClient()

			o := fake.NewSimpleClusterOverview()
			manager := modulefake.NewStubManager("default", []module.Module{o})

			d, err := newDash(listener, namespace, uiURL, nsClient, manager)
			require.NoError(t, err)

			d.willOpenBrowser = false
			d.defaultHandler = defaultHandler

			var runErr error
			ch := make(chan bool, 1)

			go func() {
				runErr = d.Run(ctx)
				ch <- true
			}()

			dashURL := fmt.Sprintf("http://%s", listener.Addr())

			resp, err := http.Get(dashURL)
			require.NoError(t, err)

			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, string(data))

			cancel()
			<-ch
			assert.NoError(t, runErr)
		})
	}
}

func TestDash_routes(t *testing.T) {
	cases := []struct {
		path         string
		expectedCode int
		expectedBody string
	}{
		{
			path:         "/",
			expectedCode: http.StatusOK,
			expectedBody: "body",
		},
		{
			path:         "/nested",
			expectedCode: http.StatusOK,
			expectedBody: "body",
		},
		{
			path:         "/api/v1/namespaces",
			expectedCode: http.StatusOK,
			expectedBody: "{}",
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("GET: %s", tc.path)
		t.Run(name, func(t *testing.T) {
			namespace := "default"
			uiURL := ""
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)

			nsClient := fake.NewNamespaceClient()

			o := fake.NewSimpleClusterOverview()
			manager := modulefake.NewStubManager("default", []module.Module{o})

			d, err := newDash(listener, namespace, uiURL, nsClient, manager)
			require.NoError(t, err)

			service := api.New(apiPathPrefix, nsClient, manager)
			d.apiHandler = service

			d.defaultHandler = func() (handler http.Handler, err error) {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "body")
				}), nil
			}

			handler, err := d.handler()
			require.NoError(t, err)

			ts := httptest.NewServer(handler)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			u.Path = tc.path

			res, err := http.Get(u.String())
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)
			data, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedBody, string(data))

		})
	}
}
