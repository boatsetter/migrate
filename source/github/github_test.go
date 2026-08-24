package github

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	st "github.com/golang-migrate/migrate/v4/source/testing"
	gh "github.com/google/go-github/v89/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var GithubTestSecret = "" // username:token

func init() {
	secrets, err := os.ReadFile(".github_test_secrets")
	if err == nil {
		GithubTestSecret = string(bytes.TrimSpace(secrets)[:])
	}
}

func TestAuthenticatedClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer secret", r.Header.Get("Authorization"))
		assert.Equal(t, "/api/v3/repos/owner/repo/contents/migrations", r.URL.Path)
		assert.Equal(t, "main", r.URL.Query().Get("ref"))
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("[]"))
		require.NoError(t, err)
	}))
	defer ts.Close()

	originalNewClient := newGithubClient
	newGithubClient = func(options ...gh.ClientOptionsFunc) (*gh.Client, error) {
		return gh.NewClient(append(options, gh.WithEnterpriseURLs(ts.URL, ts.URL))...)
	}
	t.Cleanup(func() { newGithubClient = originalNewClient })

	_, err := (&Github{}).Open("github://x-access-token:secret@owner/repo/migrations#main")
	require.NoError(t, err)
}

func TestEmptyAccessToken(t *testing.T) {
	_, err := (&Github{}).Open("github://x-access-token:@owner/repo/migrations#main")
	require.Error(t, err)
}

func TestAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer ts.Close()

	originalNewClient := newGithubClient
	newGithubClient = func(options ...gh.ClientOptionsFunc) (*gh.Client, error) {
		return gh.NewClient(append(options, gh.WithEnterpriseURLs(ts.URL, ts.URL))...)
	}
	t.Cleanup(func() { newGithubClient = originalNewClient })

	_, err := (&Github{}).Open("github://x-access-token:secret@owner/repo/migrations#main")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "401")
}

func Test(t *testing.T) {
	if len(GithubTestSecret) == 0 {
		t.Skip("test requires .github_test_secrets")
	}

	g := &Github{}
	d, err := g.Open("github://" + GithubTestSecret + "@mattes/migrate_test_tmp/test#452b8003e7")
	if err != nil {
		t.Fatal(err)
	}

	st.Test(t, d)
}

func TestDefaultClient(t *testing.T) {
	g := &Github{}
	owner := "golang-migrate"
	repo := "migrate"
	path := "source/github/examples/migrations"

	url := fmt.Sprintf("github://%s/%s/%s", owner, repo, path)
	d, err := g.Open(url)
	if err != nil {
		t.Fatal(err)
	}

	ver, err := d.First()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, uint(1085649617), ver)

	ver, err = d.Next(ver)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, uint(1185749658), ver)
}
