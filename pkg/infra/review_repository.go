package infra

import (
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"time"

	"github.com/google/go-github/github"
	"github.com/ningenMe/mami-interface/nina-api-grpc/mami"
	"github.com/ningenme/nina-api/pkg/domainmodel" //TODO domainmodelをリポジトリに切り出す
)

type ReviewRepository struct{}

const NinaApiHost = "nina-api.ningenme.net:443"

// TODO 責務が大きいので分割
func (ReviewRepository) GetContributionList(client *github.Client, ctx context.Context, pullRequest *github.PullRequest, startTime time.Time, endTime time.Time) []*domainmodel.Contribution {

	var contributionList []*domainmodel.Contribution

	org := pullRequest.GetHead().GetRepo().GetOwner().GetLogin()
	repo := pullRequest.GetHead().GetRepo().GetName()
	number := pullRequest.GetNumber()

	//prも追加
	{
		contribution := &domainmodel.Contribution{
			ContributedAt: pullRequest.GetCreatedAt(),
			Organization:  org,
			Repository:    repo,
			User:          pullRequest.GetUser().GetLogin(),
			Status:        "CREATED_PULL_REQUEST",
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

			contribution := &domainmodel.Contribution{
				ContributedAt: review.GetSubmittedAt(),
				Organization:  org,
				Repository:    repo,
				User:          review.GetUser().GetLogin(),
				Status:        review.GetState(),
			}
			contributionList = append(contributionList, contribution)
		}

		time.Sleep(1 * time.Second)

		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}

	return contributionList
}

func (ReviewRepository) PostContributionList(ctx context.Context, contributionList []*domainmodel.Contribution) {
	cc, err := grpc.Dial(
		NinaApiHost,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})),
		grpc.WithBlock(),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cc.Close()

	client := mami.NewGithubContributionServiceClient(cc)
	stream, err := client.Post(ctx)

	if err != nil {
		fmt.Println(err)
		return
	}

	partitionedList := domainmodel.PartitionedList[domainmodel.Contribution](contributionList, 20)
	for idx, list := range partitionedList {
		tmpList := []*mami.Contribution{}
		for _, co := range list {
			tmpList = append(tmpList, &mami.Contribution{
				ContributedAt: co.ContributedAt.Format(time.RFC3339),
				Organization: co.Organization,
				Repository: co.Repository,
				User: co.User,
				Status: co.Status,
			})
		}
		if err := stream.Send(&mami.PostGithubContributionRequest{
			Contributions: tmpList,
		}); err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		fmt.Println(idx+1 , "/" , len(partitionedList))
		time.Sleep(time.Second * 2)

	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return
	}
}

func (ReviewRepository) DeleteContributionList(ctx context.Context, startTime time.Time, endTime time.Time) {
	cc, err := grpc.Dial(
		NinaApiHost,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})),
		grpc.WithBlock(),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cc.Close()

	client := mami.NewGithubContributionServiceClient(cc)
	_, err = client.Delete(ctx, &mami.DeleteGithubContributionRequest{
		StartAt: startTime.Format(time.RFC3339),
		EndAt: endTime.Format(time.RFC3339),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
