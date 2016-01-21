package github

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Config struct {
	Token        string
	Organization string
}

type GithubClient struct {
	organization string
	client       *github.Client
}

// Client configures and returns a fully initialized GithubClient
func (c *Config) Client() (interface{}, error) {
	var client GithubClient
	client.organization = c.Organization
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client.client = github.NewClient(tc)
	return &client, nil
}
