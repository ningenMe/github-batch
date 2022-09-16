package cmd

import (
	"context"
	"fmt"
	"github.com/ningenme/nina-api/pkg/domainmodel"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/ningenme/nina-batch/pkg/infra"
)

const layout = "2006-01-02 15:04:05"
const startTimeSuffix = " 00:00:00"
const endTimeSuffix = " 23:59:59"
const LogIndent = "        "

type LoginUserContributionUsecase struct{}

func (LoginUserContributionUsecase) Execute(personalAccessToken string, startTimeString string, endTimeString string) {

	startTime, endTime := getPeriod(startTimeString, endTimeString)

	ctx := context.Background()

	userRepository := infra.UserRepository{}
	pullRequestRepository := infra.PullRequestRepository{}
	reviewRepository := infra.ReviewRepository{}

	//認証を行う
	fmt.Println("authentication start")
	client := userRepository.GetAuthenticatedClient(personalAccessToken, ctx)
	loginUserName := userRepository.GetLoginUserName(client, ctx)

	//repositoryの一覧を取得
	fmt.Println("getting repository list start")
	repositoryList := pullRequestRepository.GetRepositoryList(client, ctx)

	//pullRequestの一覧を取得
	fmt.Println("getting pullRequest list start")
	var pullRequestList []*github.PullRequest
	for _, repository := range repositoryList {
		org := repository.Owner.GetLogin()
		repo := repository.GetName()
		tmpPullRequestList := pullRequestRepository.GetPullRequestList(client, ctx, org, repo, startTime, endTime)
		pullRequestList = append(pullRequestList, tmpPullRequestList...)

		fmt.Println(LogIndent, org, repo, len(tmpPullRequestList))
	}
	//{
	//	org := "ningenMe"
	//	repo := "zeus"
	//	tmpPullRequestList := pullRequestRepository.GetPullRequestList(client, ctx, org, repo, startTime, endTime)
	//	pullRequestList = append(pullRequestList, tmpPullRequestList...)
	//}

	//contributionの一覧を取得
	fmt.Println("getting contribution list start")
	var contributionList []*domainmodel.Contribution
	for _, pullRequest := range pullRequestList {
		tmpContributionList := reviewRepository.GetContributionList(client, ctx, pullRequest, startTime, endTime)

		//loginUserNameでfilter
		for _, contribution := range tmpContributionList {
			if contribution.User != loginUserName {
				continue
			}
			contributionList = append(contributionList, contribution)
		}
	}

	//contributionを削除
	fmt.Println("deleting contribution list start")
	reviewRepository.DeleteContributionList(ctx, startTime, endTime)

	//contributionを永続化
	fmt.Println("inserting contribution list start")
	reviewRepository.PostContributionList(ctx, contributionList)

	fmt.Println(loginUserName)
	fmt.Println(len(repositoryList))
	os.Exit(0)
}

// TODO 適切な置き場所を考える
func getPeriod(startTimeString string, endTimeString string) (time.Time, time.Time) {
	location, _ := time.LoadLocation("Asia/Tokyo")
	startTime, _ := time.ParseInLocation(layout, startTimeString+startTimeSuffix, location)
	endTime, _ := time.ParseInLocation(layout, endTimeString+endTimeSuffix, location)
	fmt.Println(startTime, endTime)
	return startTime, endTime
}
