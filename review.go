package fastlane

import "encoding/json"

type Review struct {
	ID             string `json:"id"`
	ProjectID      string `json:"project_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Author         User   `json:"author"`
	WebURL         string `json:"web_url"`
	Approvals      []User `json:"approvals"`
	CanBeMerged    bool   `json:"can_be_merged"`
	Merged         bool   `json:"merged"`
	SHA            string `json:"sha"`
	MergeCommitSHA string `json:"merge_commit_sha"`
}

func (r Review) IsApproved() bool {
	return len(r.Approvals) > 0
}

func (r Review) MergeEnabled() bool {
	return r.IsApproved() && r.CanBeMerged
}

type enhancedReview Review

func (r Review) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		enhancedReview
		MergeEnabled bool `json:"merge_enabled"`
	}{
		enhancedReview: enhancedReview(r),
		MergeEnabled:   r.MergeEnabled(),
	})
}

type Pipeline struct {
	ID       string  `json:"id"`
	Project  string  `json:"project"`
	ReviewID string  `json:"review_id"`
	Running  bool    `json:"running"`
	Success  bool    `json:"success"`
	Status   string  `json:"status"`
	WebURL   string  `json:"web_url"`
	Stages   []Stage `json:"stages"`
	Coverage float64 `json:"coverage"`
}

type Stage struct {
	Name   string `json:"name"`
	Jobs   []Job  `json:"jobs"`
}

type Job struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	WebURL string `json:"web_url"`
}
