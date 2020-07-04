# PR labeler

Pull Request Labeler - A github action to add labels on pull request(use commits files)

> Project is refer from https://github.com/paulfantom/periodic-labeler

## Usage

### Create `.github/labeler.yml`

Create a `.github/labeler.yml` file with a list of labels and match globs to match to apply the label.

The key is the name of the label in your repository that you want to add (eg: "merge conflict", "needs-updating") and the value is the path (glob) of the changed files (eg: `src/**/*`, `tests/*.spec.js`)

#### Basic Examples

```yml
# Add 'label1' to any changes within 'example' folder or any subfolders
label1:
  - example/**/*
# Add 'label2' to any file changes within 'example2' folder
label2: example2/*
```

#### Common Examples

```yml
# Add 'repo' label to any root file changes
repo:
  - ./*
  
# Add '@domain/core' label to any change within the 'core' package
@domain/core:
  - package/core/*
  - package/core/**/*
# Add 'test' label to any change to *.spec.js files within the source dir
test:
  - src/**/*.spec.js
```

### Create Workflow

Create a workflow (eg: `.github/workflows/labeler.yml` see [Creating a Workflow file](https://help.github.com/en/articles/configuring-a-workflow#creating-a-workflow-file)) to utilize the labeler action with content:

```yml
name: "Pull Request Labeler"
on:
- pull_request
jobs:
  labels:
    runs-on: ubuntu-latest
    steps:
    - uses: inherelab/pr-labeler@master
      env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GITHUB_REPO: ${{ github.repository }}
          LABEL_CONFIG: .github/labeler.yml
```

_Note: This grants access to the `GITHUB_TOKEN` so the action can make calls to GitHub's rest API_

## Refer Projects

- https://github.com/actions/labeler
- https://github.com/paulfantom/periodic-labeler

## Refer Documents

- https://docs.github.com/en/actions
- https://developer.github.com/v3/pulls/#list-pull-requests-files
- https://docs.github.com/en/actions/configuring-and-managing-workflows/using-environment-variables
- https://docs.github.com/en/actions/reference/context-and-expression-syntax-for-github-actions#github-context
- https://docs.github.com/en/actions/reference/events-that-trigger-workflows#pull-request-event-pull_request
