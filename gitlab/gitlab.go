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
	trueBoolean       = true
)

type Gitlab struct {
	client *gitlab.Client
}

func New(client *gitlab.Client) *Gitlab {
	return &Gitlab{client: client}
}

func (g *Gitlab) ListOpenByAuthor(ctx context.Context, username string) ([]fastlane.Review, error) {
	opts := gitlab.ListMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: maxPerPage,
		},
		AuthorUsername:         &username,
		State:                  &openedState,
		WithMergeStatusRecheck: &trueBoolean,
	}

	var reviews []fastlane.Review

	for {
		mrs, res, err := g.client.MergeRequests.ListMergeRequests(&opts, gitlab.WithContext(ctx))
		if err != nil {
			return nil, err
		}
		if reviews == nil && res.TotalItems > 0 {
			reviews = make([]fastlane.Review, 0, res.TotalItems)
		}
		for _, mr := range mrs {
			r, err := buildReview(ctx, g.client, mr)
			if err != nil {
				return nil, err
			}
			reviews = append(reviews, r)
		}
		if res.NextPage == 0 {
			return reviews, nil
		}
		opts.Page = res.NextPage
	}
}

func (g *Gitlab) Merge(ctx context.Context, r fastlane.Review) (fastlane.Review, error) {
	opts := &gitlab.AcceptMergeRequestOptions{
		Squash: &trueBoolean,
		//SHA:                      &r.SHA,
		ShouldRemoveSourceBranch: &trueBoolean,
	}

	id, err := strconv.Atoi(r.ID)
	if err != nil || id <= 0 {
		return fastlane.Review{}, fmt.Errorf("review ID is not a valid gitlab merge request ID")
	}
	projectID, err := strconv.Atoi(r.ProjectID)
	if err != nil || projectID <= 0 {
		return fastlane.Review{}, fmt.Errorf("project ID is not a valid gitlab project ID")
	}

	m, _, err := g.client.MergeRequests.AcceptMergeRequest(projectID, id, opts, gitlab.WithContext(ctx))
	if err != nil {
		return fastlane.Review{}, err
	}

	return buildReviewWithApprovals(m, r.Approvals), nil
}

func (g *Gitlab) GetMergePipeline(ctx context.Context, r fastlane.Review) (fastlane.Pipeline, error) {
	sha := r.MergeCommitSHA
	if sha == "" {
		return fastlane.Pipeline{}, fmt.Errorf("no merge commit sha for review %v", r.ID)
	}

	projectID, err := strconv.Atoi(r.ProjectID)
	if err != nil || projectID <= 0 {
		return fastlane.Pipeline{}, fmt.Errorf("project ID is not a valid gitlab project ID")
	}

	c, _, err := g.client.Commits.GetCommit(projectID, sha, gitlab.WithContext(ctx))
	if err != nil {
		return fastlane.Pipeline{}, fmt.Errorf("canot get merge commit info: %w", err)
	}

	if c.LastPipeline == nil {
		return fastlane.Pipeline{}, fmt.Errorf("no pipelines for merge commit on review: %v", r.ID)
	}

	p, _, err := g.client.Pipelines.GetPipeline(projectID, c.LastPipeline.ID, gitlab.WithContext(ctx))
	if err != nil {
		return fastlane.Pipeline{}, fmt.Errorf("canot get pipeline %d info: %w", c.LastPipeline.ID, err)
	}

	opts := &gitlab.ListJobsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: maxPerPage,
		},
		IncludeRetried: false, // gitlab 13.9?
	}

	j, _, err := g.client.Jobs.ListPipelineJobs(projectID, c.LastPipeline.ID, opts, gitlab.WithContext(ctx))
	if err != nil {
		return fastlane.Pipeline{}, fmt.Errorf("canot get pipeline %d jobs: %w", c.LastPipeline.ID, err)
	}

	var cov float64
	if c, err := strconv.ParseFloat(p.Coverage, 32); err == nil {
		cov = c
	}

	stages, err := buildStages(j)
	if err != nil {
		return fastlane.Pipeline{}, err
	}

	return fastlane.Pipeline{
		ID:       strconv.Itoa(p.ID),
		Project:  strconv.Itoa(p.ProjectID),
		Running:  p.StartedAt != nil && p.FinishedAt == nil,
		Success:  p.Status == "success",
		WebURL:   p.WebURL,
		Coverage: cov,
		Stages:   stages,
	}, nil
}

func getMergeRequestApprovals(ctx context.Context, git *gitlab.Client, mr *gitlab.MergeRequest) ([]fastlane.User, error) {
	mra, _, err := git.MergeRequests.GetMergeRequestApprovals(mr.ProjectID, mr.IID, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var approvals []fastlane.User
	for _, u := range mra.ApprovedBy {
		approvals = append(approvals, user(u.User))
	}

	return approvals, nil
}

func user(u *gitlab.BasicUser) fastlane.User {
	return fastlane.User{
		ID:        strconv.Itoa(u.ID),
		Username:  u.Username,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
	}
}

func buildReview(ctx context.Context, git *gitlab.Client, mr *gitlab.MergeRequest) (fastlane.Review, error) {
	approvals, err := getMergeRequestApprovals(ctx, git, mr)
	if err != nil {
		return fastlane.Review{}, fmt.Errorf("cannot get MR %q approvals: %w", mr.Title, err)
	}
	return buildReviewWithApprovals(mr, approvals), nil
}

func buildReviewWithApprovals(mr *gitlab.MergeRequest, approvals []fastlane.User) fastlane.Review {
	return fastlane.Review{
		ID:             strconv.Itoa(mr.IID),
		ProjectID:      strconv.Itoa(mr.ProjectID),
		Author:         user(mr.Author),
		Title:          mr.Title,
		Description:    mr.Description,
		CanBeMerged:    mr.MergeStatus == canBeMergedStatus, // cannot_be_merged_recheck => f*** || unchecked ??
		WebURL:         mr.WebURL,
		SHA:            mr.SHA,
		MergeCommitSHA: mr.MergeCommitSHA,
		Approvals:      approvals,
	}
}

func buildStages(jobs []*gitlab.Job) (stages []fastlane.Stage, err error) {
	cache := make(map[string]int)

	// reverse jobs, older jobs first.
	for i, j := 0, len(jobs)-1; i < j; i, j = i+1, j-1 {
		jobs[i], jobs[j] = jobs[j], jobs[i]
	}

	for _, job := range jobs {
		pos, ok := cache[job.Stage]
		if !ok {
			pos = len(stages)
			cache[job.Stage] = pos
			stages = append(stages, fastlane.Stage{Name: job.Stage})
		}

		status, err := status(stages[pos].Status, job.Status)
		if err != nil {
			return nil, err
		}

		stages[pos].Jobs = append(stages[pos].Jobs, fastlane.Job{
			Name:   job.Name,
			Status: status,
			WebURL: job.WebURL,
		})
	}

	return stages, nil
}

var statusPriority = map[string]int{
	"success":  0,
	"pending":  1,
	"manual":   2,
	"skipped":  3,
	"canceled": 4,
	"running":  5,
	"failed":   6,
}

func status(current fastlane.Status, new string) (fastlane.Status, error) {
	a := statusPriority[string(current)]
	b, ok := statusPriority[new]
	if !ok {
		return "", fmt.Errorf("unknown status %q", new)
	}

	if b > a {
		return fastlane.Status(new), nil
	}

	return current, nil
}
