package fastlane

type Review struct {
	ID          string
	Title       string
	WebURL      string
	Approvals   []string
	CanBeMerged bool
}

func (r Review) IsApproved() bool {
	return len(r.Approvals) > 0
}
