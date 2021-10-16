package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/hmoragrega/fastlane/pushover"
	"github.com/hmoragrega/fastlane/ws"
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
		po     *pushover.Client
	)

	if os.Getenv("PUSHOVER_ENABLED") == "true" {
		po = pushover.New(
			os.Getenv("PUSHOVER_USER"),
			os.Getenv("PUSHOVER_TOKEN"),
			os.Getenv("PUSHOVER_BASE_URL"),
			pushover.WithSound(os.Getenv("PUSHOVER_SOUND")),
		)
	}
	/*
		p, err := git.GetMergePipeline(ctx, fastlane.Review{
			ID: "293",
			ProjectID: "196",
			MergeCommitSHA: "63a7b45552558a30c44d92f2af0abbeb7736902a",
		})
		fmt.Println(err)
		if err == nil {
			b, err := json.MarshalIndent(p, "", "  ")
			fmt.Println(string(b), err)
		}

		return
	*/
	rtr.GET("/v1/reviews", httpapi.ListOpenReviews(syncer))
	rtr.GET("/v1/stats", httpapi.GetStats(syncer))
	rtr.GET("/v1/ws", ws.Ws(syncer))

	svr := &http.Server{
		Addr:    ":3000",
		Handler: rtr,
	}

	os.Exit(serve(ctx, stop, syncer, svr, po))
}

func serve(ctx context.Context, stop func(), syncer *fastlane.Syncer, svr *http.Server, po *pushover.Client) (code int) {
	var wg sync.WaitGroup

	if po != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fastlane.PushSystemNotifications(ctx, po, syncer.Subscribe())
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer stop()

		if err := syncer.KeepUpdated(ctx, 10*time.Second); err != nil {
			log.Printf("error: %v", err)
			code ^= 1 << 0
		}
	}()

	wg.Add(1)
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
