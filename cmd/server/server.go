package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hmoragrega/fastlane"
	"github.com/hmoragrega/fastlane/gitlab"
	gitlabsdk "github.com/hmoragrega/go-gitlab"
)

func main() {
	ctx := context.Background()
	git, err := buildGit()
	requireNoError(err)

	mrs, err := git.ListOpenByAuthor(ctx, "hilari.jimenez")
	requireNoError(err)

	b, err := json.MarshalIndent(mrs, "", "  ")
	requireNoError(err)

	fmt.Println(string(b))
}

func buildGit() (fastlane.Git, error) {
	var opts []gitlabsdk.ClientOptionFunc
	if u := os.Getenv("GITLAB_BASE_URL"); u != "" {
		opts = append(opts, gitlabsdk.WithBaseURL(u))
	}

	c, err := gitlabsdk.NewClient(os.Getenv("GITLAB_ACCESS_TOKEN"), opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot build gitlab SDK client")
	}

	return gitlab.New(c), nil
}

func requireNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
