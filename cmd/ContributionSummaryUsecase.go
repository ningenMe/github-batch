package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/ningenme/nina/pkg/domainservice"
)

func Execute(personalAccessToken string) {

	//認証を行う
	client := domainservice.GetAuthenticatedClient(personalAccessToken)

	// repositoryList, _, err := client.Repositories.List(context.TODO(), "", nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }
	// for _, repository := range repositoryList {
	// 	fmt.Println(repository.GetURL())
	// 	organization := repository.Owner.GetLogin()
	// 	repositoryName := repository.GetName()
	// 	fmt.Println(organization)
	// 	fmt.Println(repositoryName)
	// }
	user, _, err := client.Users.Get(context.TODO(), "")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	loginUserName := user.GetLogin()
	fmt.Println(loginUserName)

	{
		org := "ningenMe"
		repo := "zeus"
		// pullRequestListOpt := &github.PullRequestListOptions{
		// 	State:       "all",
		// 	ListOptions: github.ListOptions{PerPage: 30},
		// }

		// for {
		// 	pullRequestList, response, err := client.PullRequests.List(context.TODO(), "ningenMe", "zeus", opt)
		// 	if err != nil {
		// 		fmt.Println(err.Error())
		// 		os.Exit(1)
		// 	}
		// 	for _, pullRequest := range pullRequestList {
		// 		fmt.Println(pullRequest.GetURL())
		// 	}

		// 	if response.NextPage == 0 {
		// 		break
		// 	}
		// 	opt.Page = response.NextPage
		// }
		// pullRequestList, _, err := client.PullRequests.List(context.TODO(), org, repo, pullRequestListOpt)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	os.Exit(1)
		// }

		// for _, pullRequest := range pullRequestList {
		// 	fmt.Println(pullRequest.GetURL())
		// 	reviewList, _, err := client.PullRequests.ListReviews(context.TODO(), org, repo, pullRequest.GetNumber(), reviewListOpt)
		// 	if err != nil {
		// 		fmt.Println(err.Error())
		// 		os.Exit(1)
		// 	}
		// 	for _, review := range reviewList {
		// 		fmt.Println(review)
		// 	}
		// }

		reviewCountMap := make(map[string]int)
		reviewListOpt := &github.ListOptions{PerPage: 30}
		for {
			reviewList, response, err := client.PullRequests.ListReviews(context.TODO(), org, repo, 1, reviewListOpt)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			for _, review := range reviewList {
				if review.GetUser().GetLogin() == loginUserName {
					if _, exist := reviewCountMap[review.GetState()]; !exist {
						reviewCountMap[review.GetState()] = 0
					}
					reviewCountMap[review.GetState()] += 1
				}
			}
			if response.NextPage == 0 {
				break
			}

			reviewListOpt.Page = response.NextPage
		}
		fmt.Println(reviewCountMap)
	}

	//ユーザ自身のリポジトリを取得
	//引数リポジトリとマージ
	//リポジトリの一覧のループ
	{
		//### 日付の一覧のループ
		{
			//#### prの一覧を取得
		}
	}
	//prの一覧のループ
	{
		//#### approveを取得
		//#### commentを取得
		//#### prの数を取得
	}
	//永続化
	//{repo,date,pr_count,review_comment_count,review_approve_count}
}
