package infra

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
)

type PullRequestRepository struct{}

func (PullRequestRepository) GetRepositoryList(client *github.Client, ctx context.Context) []*github.Repository {
	opt := &github.RepositoryListOptions{
		Affiliation: "owner,collaborator",
		ListOptions: github.ListOptions{PerPage: 30},
	}
	var repositoryList []*github.Repository
	for {
		tmpRepositoryList, response, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		for _, repository := range tmpRepositoryList {
			repositoryList = append(repositoryList, repository)
			fmt.Println(repository.Owner.GetLogin(), repository.GetName())
		}
		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}
	return repositoryList
}

func (PullRequestRepository) GetPullRequestList(client *github.Client, ctx context.Context, repository *github.Repository) []*github.PullRequest {

	org := repository.Owner.GetLogin()
	repo := repository.GetName()
	opt := &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 30},
	}

	var pullRequestList []*github.PullRequest
	for {
		tmpPullRequestList, response, err := client.PullRequests.List(ctx, org, repo, opt)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		pullRequestList = append(pullRequestList, tmpPullRequestList...)
		time.Sleep(1 * time.Second)

		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}

	return pullRequestList
}
