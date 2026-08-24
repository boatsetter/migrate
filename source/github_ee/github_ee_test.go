package github_ee

import (
	"net/http"
	"net/http/httptest"
	nurl "net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !assert.True(t, ok) || !assert.Equal(t, "foo", username) || !assert.Equal(t, "bar", password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.URL.Path != "/api/v3/repos/mattes/migrate_test_tmp/contents/test" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if ref := r.URL.Query().Get("ref"); ref != "452b8003e7" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte("[]"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))
	defer ts.Close()

	u, err := nurl.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	g := &GithubEE{}
	_, err = g.Open("github-ee://foo:bar@" + u.Host + "/mattes/migrate_test_tmp/test?verify-tls=false#452b8003e7")

	if err != nil {
		t.Fatal(err)
	}
}

func TestEnterpriseURLs(t *testing.T) {
	g := &GithubEE{}
	client, err := g.createGithubClient("github.example.com", "foo", "bar", true)
	require.NoError(t, err)
	assert.Equal(t, "https://github.example.com/api/v3/", client.BaseURL())
	assert.Equal(t, "https://uploads.github.example.com/", client.UploadURL())
}

func TestAPIError(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer ts.Close()

	u, err := nurl.Parse(ts.URL)
	require.NoError(t, err)

	_, err = (&GithubEE{}).Open("github-ee://foo:bar@" + u.Host + "/owner/repo/migrations?verify-tls=false#main")
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "401"), err.Error())
}
