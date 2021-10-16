package fastlane

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Stats struct {
	UpdateCount   uint64     `json:"update_count"`
	LastUpdated   *time.Time `json:"last_updated"`
	Running       bool       `json:"running"`
	Subscriptions int        `json:"subscriptions"`
}

type Syncer struct {
	git    Git
	author string

	stats         Stats
	openReviews   []Review
	mergedReviews []ReviewMerged
	subscribers   []chan Event
	merged        chan ReviewMerged
	mx            sync.RWMutex
}

func NewSync(git Git, author string) *Syncer {
	return &Syncer{
		author: author,
		git:    git,
	}
}

func (s *Syncer) Subscribe() <-chan Event {
	s.mx.Lock()
	defer s.mx.Unlock()

	c := make(chan Event)

	// send the current reviews right away.
	go func() {
		c <- Event{Name: ReviewsEventName, Data: s.openReviews}
		c <- Event{Name: ReviewsMergedEventName, Data: s.mergedReviews}
	}()

	s.subscribers = append(s.subscribers, c)
	s.stats.Subscriptions = len(s.subscribers)

	return c
}

func (s *Syncer) Unsubscribe(c <-chan Event) {
	s.mx.Lock()
	defer s.mx.Unlock()

	for i, sub := range s.subscribers {
		if sub == c {
			close(s.subscribers[i])
			s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
			s.stats.Subscriptions = len(s.subscribers)
			break
		}
	}
}

func (s *Syncer) KeepUpdated(ctx context.Context, frequency time.Duration) error {
	if err := s.start(); err != nil {
		return err
	}

	defer s.closeSubscribers()

	t := time.NewTimer(frequency)
	t.Reset(0)
	defer t.Stop()

	// TODO remove, test pipeline
	s.mergedReviews = append(s.mergedReviews, ReviewMerged{Review: Review{
		Title:          "Fastlane test n.002",
		ID:             "293",
		Description:    "Partners can have whitelisted IPs optionally It accepts both normal IPv4 and v6 or in CIDR notation Next: Part 2 - Be able to Update the IPs as Admin (or even the partner?)",
		ProjectID:      "196",
		MergeCommitSHA: "63a7b45552558a30c44d92f2af0abbeb7736902a",
	}})

	for {
		select {
		case <-ctx.Done():
			return nil
		case rm := <-s.merged:
			// todo create method
			s.mx.Lock()
			s.mergedReviews = append([]ReviewMerged{rm}, s.mergedReviews...)
			s.mx.Unlock()
		case <-t.C:
			t.Reset(frequency)

			events, err := s.pipelines(ctx)
			if errors.Is(err, context.Canceled) {
				return nil
			}
			if err != nil {
				return err
			}
			s.pushUpdates(ctx, events...)

			events, err = s.update(ctx)
			if errors.Is(err, context.Canceled) {
				return nil
			}
			if err != nil {
				return err
			}
			s.pushUpdates(ctx, events...)
		}
	}
}

func (s *Syncer) update(ctx context.Context) ([]Event, error) {
	mrs, err := s.git.ListOpenByAuthor(ctx, s.author)
	if err != nil {
		return nil, err
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	events := []Event{{Name: ReviewsEventName, Data: mrs}}
	events = append(events, calculateReviewEvents(s.openReviews, mrs)...)

	now := time.Now()
	s.stats.LastUpdated = &now
	s.stats.UpdateCount++
	s.openReviews = mrs

	return events, nil
}

func (s *Syncer) OpenReviews() []Review {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.openReviews
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
		return fmt.Errorf("alredy running")
	}

	s.stats.Running = true
	s.merged = make(chan ReviewMerged)
	return nil
}

func (s *Syncer) pushUpdates(ctx context.Context, updates ...Event) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	var wg sync.WaitGroup
	for _, sub := range s.subscribers {

		wg.Add(1)
		go func(sub chan<- Event) {
			defer wg.Done()
			for _, e := range updates {
				select {
				case <-ctx.Done():
				case sub <- e:
				}
			}
		}(sub)
	}

	wg.Wait()
}

func (s *Syncer) closeSubscribers() {
	s.mx.RLock()
	defer s.mx.RUnlock()
	for _, sub := range s.subscribers {
		close(sub)
	}
}

// Handle processes an event, it may return
// other events as response.
func (s *Syncer) Handle(ctx context.Context, evt Event) []Event {
	switch name := evt.Name; name {
	case MergeEventName:
		return s.merge(ctx, evt.Data.(string))

	default:
		return []Event{notificationEvent(Notification{
			Message: fmt.Sprintf("unknown event name %q", name),
			Type:    WarningType,
		})}
	}
}

func (s *Syncer) merge(ctx context.Context, id string) []Event {
	r := s.getOpenReview(id)
	if r == nil {
		return []Event{notificationEvent(Notification{
			Message: fmt.Sprintf("review %v not found", id),
			Type:    WarningType,
		})}
	}

	merged, err := s.git.Merge(ctx, *r)
	if err != nil {
		return []Event{notificationEvent(Notification{
			Message: fmt.Sprintf("review %v could not be merged: %v", id, err),
			Type:    ErrorType,
		})}
	}

	go func() {
		// remove the review and push the new list.
		s.removeReview(id)
		s.pushUpdates(ctx,
			Event{Name: ReviewsEventName, Data: s.getOpenReviews()},
			Event{Name: ReviewMergedEventName, Data: ReviewMerged{Review: merged}},
		)
	}()

	s.merged <- ReviewMerged{Review: merged}

	return []Event{notificationEvent(Notification{
		Message: fmt.Sprintf("review %v has been merged", id),
		Type:    SuccessType,
	})}
}

func (s *Syncer) getOpenReview(id string) *Review {
	s.mx.RLock()
	defer s.mx.RUnlock()

	for _, r := range s.openReviews {
		if r.ID == id {
			return &r
		}
	}

	return nil
}

func (s *Syncer) getOpenReviews() []Review {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.openReviews
}

func (s *Syncer) removeReview(id string) {
	s.mx.Lock()
	defer s.mx.Unlock()

	for i, r := range s.openReviews {
		if r.ID == id {
			s.openReviews = append(s.openReviews[:i], s.openReviews[i+1:]...)
			break
		}
	}
}

func (s *Syncer) pipelines(ctx context.Context) ([]Event, error) {
	type result struct {
		review   Review
		pipeline Pipeline
		index    int
		err      error
	}

	s.mx.RLock()
	current := s.mergedReviews
	s.mx.RUnlock()

	results := make(chan result, len(current))

	var wg sync.WaitGroup
	for i, r := range current {
		wg.Add(1)
		go func(i int, r ReviewMerged) {
			defer wg.Done()
			p, err := s.git.GetMergePipeline(ctx, r.Review)
			results <- result{
				review:   r.Review,
				pipeline: p,
				index:    i,
				err:      err,
			}
		}(i, r)
	}

	events := make([]Event, len(current))
	go func() {
		for res := range results {
			if res.err != nil {
				events = append(events, notificationEvent(Notification{
					Message: fmt.Sprintf("cannot get pipeline for merged review %s: %v", res.review.ID, res.err),
					Type:    ErrorType,
				}))
				continue
			}
			current[res.index] = ReviewMerged{
				Review:      res.review,
				Pipeline:    res.pipeline,
				HasPipeline: true,
			}
		}
	}()

	wg.Wait()

	s.mx.Lock()
	s.mergedReviews = current
	s.mx.Unlock()

	events = append(events, Event{Name: ReviewsMergedEventName, Data: current})

	return events, nil
}

func notificationEvent(n Notification) Event {
	n.ID = uuid.New()
	return Event{
		Name: NotificationEventName,
		Data: n,
	}
}
