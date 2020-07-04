package main

import (
	"context"
	"os"
	"strings"

	"github.com/gobwas/glob"
	"github.com/google/go-github/v32/github"
	"github.com/inherelab/pr-labeler/slog"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func getCurrentLabels(pr *github.PullRequest) []string {
	var labelSet []string
	for _, l := range pr.Labels {
		labelSet = append(labelSet, l.GetName())
	}
	return labelSet
}

func containsLabels(expected []string, current []string) bool {
	for _, e := range expected {
		if !contains(current, e) {
			return false
		}
	}
	return true
}

// Get files and labels matchers, output labels
func matchFiles(labelsMatch map[string][]glob.Glob, files []*github.CommitFile) []string {
	var labelSet []string
	set := make(map[string]bool)
	for _, file := range files {
		for label, matchers := range labelsMatch {
			if set[label] {
				continue
			}
			for _, m := range matchers {
				if m.Match(file.GetFilename()) {
					set[label] = true
					labelSet = append(labelSet, label)
					break
				}
			}
		}
	}
	return labelSet
}

func buildLabelMatchers(from string) (map[string][]glob.Glob, error) {
	var config map[string][]string
	if err := yaml.Unmarshal([]byte(from), &config); err != nil {
		return nil, err
	}

	matchers := make(map[string][]glob.Glob, len(config))

	for label, patterns := range config {
		for _, p := range patterns {
			m, err := glob.Compile(p, '/')
			if err != nil {
				return nil, err
			}
			matchers[label] = append(matchers[label], m)
		}
	}

	return matchers, nil
}

func main() {
	var owner, repo, ghToken string

	repoSlug := os.Getenv("GITHUB_REPOSITORY")
	ghToken = os.Getenv("GITHUB_TOKEN")

	confPath, exists := os.LookupEnv("LABEL_MAPPINGS_FILE")
	if !exists {
		confPath = ".github/labeler.yml"
	}

	s := strings.Split(repoSlug, "/")
	owner = s[0]
	repo = s[1]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		// TODO: access token should be passed as CLI parameter
		&oauth2.Token{AccessToken: ghToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	content, _, _, err := client.Repositories.GetContents(context.Background(), owner, repo, confPath, nil)
	if err != nil {
		slog.Fatal(err)
	}

	yamlFile, err := content.GetContent()
	if err != nil {
		slog.Fatal(err)
	}

	labelMatchers, err := buildLabelMatchers(yamlFile)
	if err != nil {
		slog.Fatal(err)
	}

	opt := &github.PullRequestListOptions{State: "open", Sort: "updated"}
	// get all pages of results
	for {
		pulls, resp, err := client.PullRequests.List(context.Background(), owner, repo, opt)
		if err != nil {
			slog.Fatal(err)
		}

		for _, pull := range pulls {
			files, _, err := client.PullRequests.ListFiles(context.Background(), owner, repo, pull.GetNumber(), nil)
			if err != nil {
				slog.Error(err)
			}

			expectedLabels := matchFiles(labelMatchers, files)
			if !containsLabels(expectedLabels, getCurrentLabels(pull)) {
				slog.Infof("PR %s/%s#%d should have following labels: %v", owner, repo, pull.GetNumber(), expectedLabels)

				_, _, err = client.Issues.AddLabelsToIssue(context.Background(), owner, repo, pull.GetNumber(), expectedLabels)
				if err != nil {
					slog.Error(err)
				}
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}
}

