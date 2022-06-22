// Copyright (C) 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package agent

import (
	"bytes"
	"context"
	"log"
	"strings"

	atlas "github.com/explore-dev/atlas-common/go/api/services"
	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/host-event-handler/handler"
	reviewpad_premium "github.com/reviewpad/reviewpad-premium/v2"
	"github.com/reviewpad/reviewpad/v2"
	"github.com/reviewpad/reviewpad/v2/collector"
	"github.com/reviewpad/reviewpad/v2/engine"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

var MixpanelToken string

const ReviewpadFile string = "reviewpad.yml"

type Env struct {
	RepoOwner        string
	RepoName         string
	Token            string
	PRNumber         int
	SemanticEndpoint string
	EventPayload     interface{}
}

func runReviewpad(prNum int, e *handler.ActionEvent, semanticEndpoint *string) {
	repo := *e.Repository
	splittedRepo := strings.Split(repo, "/")
	repoOwner := splittedRepo[0]
	repoName := splittedRepo[1]
	eventPayload, err := github.ParseWebHook(*e.EventName, *e.EventPayload)

	if err != nil {
		log.Print(err)
		return
	}

	env := &Env{
		RepoOwner:        repoOwner,
		RepoName:         repoName,
		Token:            *e.Token,
		PRNumber:         prNum,
		SemanticEndpoint: *semanticEndpoint,
		EventPayload:     eventPayload,
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: env.Token})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	clientGQL := githubv4.NewClient(tc)

	pullRequest, _, err := client.PullRequests.Get(ctx, env.RepoOwner, env.RepoName, env.PRNumber)
	if err != nil {
		log.Print(err)
		return
	}

	if pullRequest.Merged != nil && *pullRequest.Merged {
		log.Print("skip execution for merged pull requests")
		return
	}

	// TODO: Extend logic to choose between base or head
	// TODO: Check for nils
	headRepoOwner := *pullRequest.Head.Repo.Owner.Login
	headRepoName := *pullRequest.Head.Repo.Name
	headRef := *pullRequest.Head.Ref

	ioReader, _, err := client.Repositories.DownloadContents(ctx, headRepoOwner, headRepoName, ReviewpadFile, &github.RepositoryContentGetOptions{
		Ref: headRef,
	})
	if err != nil {
		log.Print(err.Error())
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(ioReader)

	file, err := reviewpad.Load(buf)
	if err != nil {
		log.Print(err.Error())
		return
	}

	collectorClient := collector.NewCollector(MixpanelToken, headRepoOwner)

	defaultOptions := grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(419430400))
	semanticConnection, err := grpc.Dial(env.SemanticEndpoint, grpc.WithInsecure(), defaultOptions)
	if err != nil {
		log.Printf("failed to dial semantic service: %v", err)
		return
	}
	defer semanticConnection.Close()
	semanticClient := atlas.NewSemanticClient(semanticConnection)

	switch file.Edition {
	case engine.PROFESSIONAL_EDITION:
		err = reviewpad_premium.Run(ctx, client, clientGQL, collectorClient, semanticClient, pullRequest, eventPayload, file, false)
	default:
		_, err = reviewpad.Run(ctx, client, clientGQL, collectorClient, pullRequest, eventPayload, file, false)
	}

	if err != nil {
		log.Print(err.Error())
		return
	}
}

func RunAction(semanticEndpoint *string, rawEvent *string) {
	event, err := handler.ParseEvent(*rawEvent)
	if err != nil {
		log.Print(err)
		return
	}

	prs := handler.ProcessEvent(event)

	for _, pr := range prs {
		runReviewpad(pr, event, semanticEndpoint)
	}
}
