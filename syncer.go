package fastlane

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Stats struct {
	UpdateCount uint64     `json:"update_count"`
	LastUpdated *time.Time `json:"last_updated"`
	Running     bool       `json:"running"`
}

type Syncer struct {
	git    Git
	author string

	stats       Stats
	openReviews []Review
	mx          sync.RWMutex
}

func NewSync(git Git, author string) *Syncer {
	return &Syncer{
		author: author,
		git:    git,
	}
}

func (s *Syncer) KeepUpdated(ctx context.Context, frequency time.Duration) error {
	if err := s.start(); err != nil {
		return err
	}

	if err := s.Update(ctx); err != nil {
		return err
	}

	t := time.NewTicker(frequency)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			err := s.Update(ctx)
			if errors.Is(err, context.Canceled) {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}

func (s *Syncer) Update(ctx context.Context) error {
	mrs, err := s.git.ListOpenByAuthor(ctx, s.author)
	if err != nil {
		return err
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	now := time.Now()
	s.openReviews = mrs
	s.stats.LastUpdated = &now
	s.stats.UpdateCount++

	return nil
}

func (s *Syncer) ListOpenReviews(_ context.Context) ([]Review, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.openReviews, nil
}

func (s *Syncer) Stats() Stats {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.stats
}

func (s *Syncer) start() error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.stats.Running {
		s.mx.Unlock()
		return fmt.Errorf("alredy syncing")
	}

	s.stats.Running = true
	return nil
}
