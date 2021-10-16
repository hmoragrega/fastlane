package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/hmoragrega/fastlane"
	"github.com/hmoragrega/fastlane/gitlab"
	"github.com/hmoragrega/fastlane/httpapi"
	gitlabsdk "github.com/hmoragrega/go-gitlab"
	"github.com/julienschmidt/httprouter"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var (
		git    = buildGit()
		syncer = fastlane.NewSync(git, os.Getenv("FASTLANE_AUTHOR"))
		rtr    = httprouter.New()
	)

	rtr.GET("/v1/reviews", httpapi.ListOpenReviews(syncer))
	rtr.GET("/v1/stats", httpapi.GetStats(syncer))

	svr := &http.Server{
		Addr:    ":8080",
		Handler: rtr,
	}

	os.Exit(serve(ctx, stop, syncer, svr))
}

func serve(ctx context.Context, stop func(), syncer *fastlane.Syncer, svr *http.Server) (code int) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer stop()

		if err := syncer.KeepUpdated(ctx, 2*time.Second); err != nil {
			log.Printf("error: %v", err)
			code ^= 1 << 0
		}
	}()

	go func() {
		defer wg.Done()
		defer stop()

		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error: %v", err)
			code ^= 1 << 1
		}
	}()

	log.Println("ready!")
	<-ctx.Done()

	if err := svr.Shutdown(context.Background()); err != nil {
		code ^= 1 << 2
	}

	wg.Wait()
	return code
}

func buildGit() fastlane.Git {
	var opts []gitlabsdk.ClientOptionFunc
	if u := os.Getenv("GITLAB_BASE_URL"); u != "" {
		opts = append(opts, gitlabsdk.WithBaseURL(u))
	}

	c, err := gitlabsdk.NewClient(os.Getenv("GITLAB_ACCESS_TOKEN"), opts...)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot build gitlab SDK client: %v", err))
	}

	return gitlab.New(c)
}
