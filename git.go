package fastlane

import (
	"context"
)

type Git interface {
	// ListOpenByAuthor list the author's open reviews.
	ListOpenByAuthor(ctx context.Context, username string) ([]Review, error)

	// Merge merges the review and return the update review.
	Merge(ctx context.Context, r Review) (Review, error)

	// GetMergePipeline returns the pipeline for a merged review.
	GetMergePipeline(ctx context.Context, r Review) (Pipeline, error)
}
