package main

import (
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/gobwas/glob"
	"github.com/google/go-github/v32/github"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/slog"
	"gopkg.in/yaml.v3"
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

// version 1
func main() {
	ghToken := os.Getenv("GITHUB_TOKEN")
	if ghToken == "" {
		slog.Fatal("please set env: GITHUB_TOKEN")
	}

	// Inner ENV see: https://docs.github.com/en/actions/configuring-and-managing-workflows/using-environment-variables
	// GITHUB_REF
	// - on PR: "refs/pull/:prNumber/merge"
	// - on push: "refs/heads/master"
	// - on push tag: "refs/tags/v0.0.1"
	ghRefer := os.Getenv("GITHUB_REF")
	// ghRefer := "refs/pull/34/merge"
	prNumber := getPrNumber(ghRefer)
	if prNumber == 0 {
		slog.Fatalf("parse PR number failed, GITHUB_REF: %s", ghRefer)
	}

	confPath := os.Getenv("LABEL_CONFIG")
	if confPath == "" {
		confPath = ".github/labeler.yml"
	}

	slog.Infof("use label config: %s", confPath)

	repoSlug := os.Getenv("GITHUB_REPOSITORY")
	nameNodes := strings.Split(repoSlug, "/")
	owner, repo := nameNodes[0], nameNodes[1]

	// create client
	gha := NewGithub(owner, repo, ghToken)

	// - fetch config file contents
	// GITHUB_SHA
	fetchContentOpts := &github.RepositoryContentGetOptions{Ref: os.Getenv("GITHUB_SHA")}
	content, _, _, err := gha.Repositories().GetContents(context.Background(), owner, repo, confPath, fetchContentOpts)
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

	// opt := &github.PullRequestListOptions{State: "open", Sort: "updated"}
	pull, _, err := gha.GetPullRequest(prNumber)
	if err != nil {
		slog.Fatal(err)
	}

	// get all changed files of PR
	files, _, err := gha.ListPullRequestFiles(prNumber, nil)
	if err != nil {
		slog.Fatal(err)
	}

	expectedLabels := matchFiles(labelMatchers, files)
	if !containsLabels(expectedLabels, getCurrentLabels(pull)) {
		slog.Infof("PR %s/%s#%d should have following labels: %v", owner, repo, prNumber, expectedLabels)

		_, _, err = gha.AddLabelsToIssue(prNumber, expectedLabels)
		if err != nil {
			slog.Error(err)
		}
	} else {
		slog.Infof("not add any labels")
	}
}

//
// match files
//

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

//
// match labels
//

//
// helper func
//

func getPrNumber(ghRefer string) int {
	// "refs/pull/:prNumber/merge"
	rg := regexp.MustCompile(`^refs/pull/(\d+)/merge`)
	ns := rg.FindStringSubmatch(ghRefer)

	if len(ns) < 2 {
		return 0
	}

	return mathutil.MustInt(ns[1])
}
