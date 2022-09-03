package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/ningenme/nina/pkg/infra"
)

func Execute(personalAccessToken string) {

	ctx := context.Background()

	userRepository := infra.UserRepository{}
	pullRequestRepository := infra.PullRequestRepository{}
	reviewRepository := infra.ReviewRepository{}

	//認証を行う
	client := userRepository.GetAuthenticatedClient(personalAccessToken, ctx)
	//ユーザ名を取得
	loginUserName := userRepository.GetLoginUserName(client, ctx)
	//repisotiryの一覧を取得
	repositoryList := pullRequestRepository.GetRepositoryList(client, ctx)
	//pullRequestの一覧を取得
	var pullRequestList []*github.PullRequest
	for _, repository := range repositoryList {
		org := repository.Owner.GetLogin()
		repo := repository.GetName()
		tmpPullRequestList := pullRequestRepository.GetPullRequestList(client, ctx, org, repo)
		pullRequestList = append(pullRequestList, tmpPullRequestList...)
	}
	//contributionの一覧を取得
	var contributionList []infra.Contribution
	for _, pullRequest := range pullRequestList {
		tmpContributionList := reviewRepository.GetContributionList(client, ctx, pullRequest)
		contributionList = append(contributionList, tmpContributionList...)
	}
	fmt.Println(len(repositoryList))
	fmt.Println(len(pullRequestList))
	fmt.Println(len(contributionList))
	fmt.Println(loginUserName)
	os.Exit(1)
}
