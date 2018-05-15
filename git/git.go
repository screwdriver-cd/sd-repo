package git

import (
	"fmt"
	"regexp"
	"strings"
)

type GitUrl struct {
	Host   string
	Org    string
	Repo   string
	Path   string
	Branch string
}

// GetCloneInfo returns the url and branch of the GitUrl
func (git *GitUrl) GetCloneInfo() (url, branch string) {
	return fmt.Sprintf("git@%+s:%+s/%+s.git", git.Host, git.Org, git.Repo), git.Branch
}

// New validates the gitUrlStr and returns a new GitUrl object
func New(gitUrlStr string) (*GitUrl, error) {
	// This would match something like git@github.com:org/repo.git/path#branch
	// path and branch are optional. If not given, default values are "" and "master"
	gitUrlRegex, _ := regexp.Compile("^git@([^/:#]+):([^/:#]+)/+([^/:#]+)\\.git(/[^#]*)?(#.+)?")
	parseResult := gitUrlRegex.FindStringSubmatch(gitUrlStr)

	if parseResult == nil {
		return nil, fmt.Errorf("Not a valid git url %+s", gitUrlStr)
	}

	gitUrl := GitUrl{
		Host:   parseResult[1],
		Org:    parseResult[2],
		Repo:   parseResult[3],
		Path:   parseResult[4],
		Branch: parseResult[5],
	}

	if gitUrl.Branch == "" {
		gitUrl.Branch = "master"
	} else {
		gitUrl.Branch = strings.TrimPrefix(gitUrl.Branch, "#")
	}

	if gitUrl.Path != "" {
		gitUrl.Path = strings.TrimLeft(gitUrl.Path, "/")
	}

	return &gitUrl, nil
}