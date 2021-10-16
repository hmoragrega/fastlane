package gitlab

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hmoragrega/fastlane"
	"github.com/hmoragrega/go-gitlab"
)

const (
	maxPerPage = 100
)

var (
	openedState       = "opened"
	canBeMergedStatus = "can_be_merged"
)

type Gitlab struct {
	client *gitlab.Client
}

func New(client *gitlab.Client) *Gitlab {
	return &Gitlab{client: client}
}

func (g *Gitlab) ListOpenByAuthor(ctx context.Context, username string) ([]fastlane.Review, error) {
	return ListOpenByAuthor(ctx, g.client, username)
}

func ListOpenByAuthor(ctx context.Context, git *gitlab.Client, username string) ([]fastlane.Review, error) {
	recheck := true
	var opts gitlab.ListMergeRequestsOptions
	opts.Page = 1
	opts.PerPage = maxPerPage
	opts.AuthorUsername = &username
	opts.State = &openedState
	opts.WithMergeStatusRecheck = &recheck

	var reviews []fastlane.Review

	for {
		mrs, res, err := git.MergeRequests.ListMergeRequests(&opts, gitlab.WithContext(ctx))
		if err != nil {
			return nil, err
		}
		if reviews == nil && res.TotalItems > 0 {
			reviews = make([]fastlane.Review, 0, res.TotalItems)
		}
		for _, mr := range mrs {
			approvals, err := getMergeRequestApprovals(ctx, git, mr)
			if err != nil {
				return nil, fmt.Errorf("cannot get MR %q approvals: %w", mr.Title, err)
			}
			reviews = append(reviews, fastlane.Review{
				ID:          strconv.Itoa(mr.ID),
				Title:       mr.Title,
				CanBeMerged: mr.MergeStatus == canBeMergedStatus, // cannot_be_merged_recheck => fuck || unchecked ??
				Approvals:   approvals,
				WebURL:      mr.WebURL,

				// UserNotesCount = number of comments
				// WorkInProgress = WIP status
				// Description    = ...
			})
		}
		if res.NextPage == 0 {
			return reviews, nil
		}
		opts.Page = res.NextPage
	}
}

func getMergeRequestApprovals(ctx context.Context, git *gitlab.Client, mr *gitlab.MergeRequest) ([]string, error) {
	mra, _, err := git.MergeRequests.GetMergeRequestApprovals(mr.ProjectID, mr.IID, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var approvals []string
	for _, a := range mra.ApprovedBy {
		approvals = append(approvals, a.User.Username)
	}

	return approvals, nil
}
