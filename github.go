package main

import (
	"context"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

type GithubApi struct {
	// owner and repo name
	owner, repo string

	client *github.Client
}

// NewGithub create
func NewGithub(owner, repo, token string) *GithubApi {
	gha := &GithubApi{
		owner: owner,
		repo:  repo,
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	// create client
	gha.client = github.NewClient(tc)
	return gha
}

// AddLabelsToIssue by PR number
func (gha *GithubApi) AddLabelsToIssue(number int, labels []string) ([]*github.Label, *github.Response, error) {
	return gha.client.Issues.AddLabelsToIssue(context.Background(), gha.owner, gha.repo, number, labels)
}

// GetPullRequest by PR number
func (gha *GithubApi) GetPullRequest(prNumber int) (*github.PullRequest, *github.Response, error) {
	return gha.client.PullRequests.Get(context.Background(), gha.owner, gha.repo, prNumber)
}

// ListPullRequestFiles by PR number
func (gha *GithubApi) ListPullRequestFiles(prNumber int, options *github.ListOptions) ([]*github.CommitFile, *github.Response, error) {
	return gha.client.PullRequests.ListFiles(context.Background(), gha.owner, gha.repo, prNumber, options)
}

// Repositories service get
func (gha *GithubApi) Repositories() *github.RepositoriesService {
	return gha.client.Repositories
}

// Repo name
func (gha *GithubApi) Repo() string {
	return gha.repo
}

// Owner name
func (gha *GithubApi) Owner() string {
	return gha.owner
}

// Client get github client
func (gha *GithubApi) Client() *github.Client {
	return gha.client
}
