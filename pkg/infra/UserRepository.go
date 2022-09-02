package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type UserRepository struct{}

func (UserRepository) GetAuthenticatedClient(accessToken string, ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}

func (UserRepository) GetLoginUserName(client *github.Client, ctx context.Context) string {
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	loginUserName := user.GetLogin()
	fmt.Println("loginUserName", loginUserName)
	return loginUserName
}
