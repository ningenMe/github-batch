package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

// TODO domainmodelに移す
type Contribution struct {
	org   string
	repo  string
	time  time.Time
	year  int
	month time.Month
	day   int
	user  string
	state string
}

type ReviewRepository struct{}

func (ReviewRepository) GetContributionList(client *github.Client, ctx context.Context, pullRequest *github.PullRequest) []Contribution {

	var contributionList []Contribution
	org := pullRequest.GetHead().GetRepo().GetOwner().GetLogin()
	repo := pullRequest.GetHead().GetRepo().GetName()
	number := pullRequest.GetNumber()

	//prも追加
	{
		contribution := Contribution{
			org:   org,
			repo:  repo,
			time:  pullRequest.GetCreatedAt(),
			year:  pullRequest.GetCreatedAt().Year(),
			month: pullRequest.GetCreatedAt().Month(),
			day:   pullRequest.GetCreatedAt().Day(),
			user:  pullRequest.GetUser().GetLogin(),
			state: "CREATED_PULL_REQUEST",
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
			contribution := Contribution{
				org:   org,
				repo:  repo,
				time:  review.GetSubmittedAt(),
				year:  review.GetSubmittedAt().Year(),
				month: review.GetSubmittedAt().Month(),
				day:   review.GetSubmittedAt().Day(),
				user:  review.GetUser().GetLogin(),
				state: review.GetState(),
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
