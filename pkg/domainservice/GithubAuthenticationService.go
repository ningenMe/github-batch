package domainservice

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GetAuthenticatedClient(accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.TODO(), ts)

	client := github.NewClient(tc)
	return client
}
