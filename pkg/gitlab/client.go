package gitlab

import (
	"github.com/lshgdut/repoctl/pkg/config"
	"github.com/xanzy/go-gitlab"
)

func NewClient() (*gitlab.Client, error) {
	// mux is the HTTP request multiplexer used with the test server.
	// mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	// server := httptest.NewServer(mux)

	config, err := config.LoadRepoctlConfig()
	if err != nil {
		return nil, err
	}

	url := config.GitlabUrl
	token := config.GitlabToken

	// client is the Gitlab client being tested.
	client, _ := gitlab.NewClient(token, gitlab.WithBaseURL(url))

	return client, err
}
