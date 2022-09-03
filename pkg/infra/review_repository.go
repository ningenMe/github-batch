package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

// TODO domainmodelに移す
type Contribution struct {
	Org   string
	Repo  string
	Time  time.Time
	User  string
	State string
}

type ReviewRepository struct{}

// TODO 責務が大きいので分割
func (ReviewRepository) GetContributionList(client *github.Client, ctx context.Context, pullRequest *github.PullRequest, startTime time.Time, endTime time.Time) []Contribution {

	var contributionList []Contribution
	org := pullRequest.GetHead().GetRepo().GetOwner().GetLogin()
	repo := pullRequest.GetHead().GetRepo().GetName()
	number := pullRequest.GetNumber()

	//prも追加
	{
		contribution := Contribution{
			Org:   org,
			Repo:  repo,
			Time:  pullRequest.GetCreatedAt(),
			User:  pullRequest.GetUser().GetLogin(),
			State: "CREATED_PULL_REQUEST",
		}
		contributionList = append(contributionList, contribution)
	}

	opt := &github.ListOptions{PerPage: 30}
	for {
		reviewList, response, err := client.PullRequests.ListReviews(
			context.Background(),
			org,
			repo,
			number,
			opt,
		)
		if err != nil {
			// fmt.Println(err.Error())
			// os.Exit(1)
			break
		}

		for _, review := range reviewList {

			//TODO periodごとドメインモデルに移して処理を共通化する
			if review.GetSubmittedAt().After(endTime) {
				continue
			}
			if review.GetSubmittedAt().Before(startTime) {
				continue
			}

			contribution := Contribution{
				Org:   org,
				Repo:  repo,
				Time:  review.GetSubmittedAt(),
				User:  review.GetUser().GetLogin(),
				State: review.GetState(),
			}
			fmt.Println(contribution)
			contributionList = append(contributionList, contribution)
		}

		fmt.Println(org, repo, number, len(reviewList), len(contributionList))
		time.Sleep(1 * time.Second)

		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}

	return contributionList
}
