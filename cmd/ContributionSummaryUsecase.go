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

	//認証を行う
	client := userRepository.GetAuthenticatedClient(personalAccessToken, ctx)
	//ユーザ名を取得
	loginUserName := userRepository.GetLoginUserName(client, ctx)
	//repisotiryの一覧を取得
	repositoryList := pullRequestRepository.GetRepositoryList(client, ctx)
	//pullRequestの一覧を取得
	var pullRequestList []*github.PullRequest
	for _, repository := range repositoryList {
		tmpPullRequestList := pullRequestRepository.GetPullRequestList(client, ctx, repository)
		pullRequestList = append(pullRequestList, tmpPullRequestList...)
	}
	//loginUserでfilter
	var loginUserPullRequestList []*github.PullRequest
	for _, pullRequest := range pullRequestList {
		if pullRequest.GetUser().GetLogin() != loginUserName {
			continue
		}
		// fmt.Println(pullRequest.GetHead().GetRepo().GetOwner().GetLogin()) //org
		// fmt.Println(pullRequest.GetHead().GetRepo().GetName())             //repo
		// fmt.Println(pullRequest.GetNumber())
		// fmt.Println(pullRequest.GetCreatedAt().Date())
		// fmt.Println(pullRequest.GetUpdatedAt().Date())
		// fmt.Println(pullRequest.GetUser().GetLogin()) //PRを出したユーザー

		loginUserPullRequestList = append(loginUserPullRequestList, pullRequest)
	}
	fmt.Println(len(loginUserPullRequestList))
	os.Exit(1)
	fmt.Println(loginUserName)
	fmt.Println(loginUserPullRequestList)

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
	// {
	// 	org := "ningenMe"
	// 	repo := "zeus"

	// 	reviewCountMap := make(map[string]int)
	// 	reviewListOpt := &github.ListOptions{PerPage: 30}
	// 	for {
	// 		reviewList, response, err := client.PullRequests.ListReviews(context.TODO(), org, repo, 1, reviewListOpt)
	// 		if err != nil {
	// 			fmt.Println(err.Error())
	// 			os.Exit(1)
	// 		}
	// 		for _, review := range reviewList {
	// 			if review.GetUser().GetLogin() == loginUserName {
	// 				if _, exist := reviewCountMap[review.GetState()]; !exist {
	// 					reviewCountMap[review.GetState()] = 0
	// 				}
	// 				reviewCountMap[review.GetState()] += 1
	// 			}
	// 		}
	// 		if response.NextPage == 0 {
	// 			break
	// 		}

	// 		reviewListOpt.Page = response.NextPage
	// 	}
	// 	fmt.Println(reviewCountMap)
	// }

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
	//{org,repo,date,loginUserName,pr_count,review_comment_count,review_approve_count}
}
