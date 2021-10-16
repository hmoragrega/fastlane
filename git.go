package fastlane

import (
	"context"
)

type Git interface {
	ListOpenByAuthor(ctx context.Context, username string) ([]Review, error)
}
