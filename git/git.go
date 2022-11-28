package git

import (
	"fmt"
	"regexp"
	"strings"
)

type GitUrl struct {
	Protocol string
	Host     string
	Org      string
	Repo     string
	Path     string
	Branch   string
}

// GetCloneInfo returns the url and branch of the GitUrl
func (git *GitUrl) GetCloneInfo() (url, branch string) {
	if git.Protocol == "https" {
		return fmt.Sprintf("https://%+s/%+s/%+s.git", git.Host, git.Org, git.Repo), git.Branch
	}

	return fmt.Sprintf("git@%+s:%+s/%+s.git", git.Host, git.Org, git.Repo), git.Branch
}

// New validates the gitUrlStr and returns a new GitUrl object
func New(gitUrlStr string) (*GitUrl, error) {
	// This would match something like git@github.com:org/repo.git/path#branch
	// path and branch are optional. If not given, default values are "" and "master"
	gitUrlRegex, _ := regexp.Compile("^(git|https)(?:@|://)([^/:#]+)(?::|/)([^/:#]+)/+([^/:#]+)\\.git(/[^#]*)?(#.+)?")
	parseResult := gitUrlRegex.FindStringSubmatch(gitUrlStr)

	if parseResult == nil {
		return nil, fmt.Errorf("Not a valid git url %+s", gitUrlStr)
	}

	gitUrl := GitUrl{
		Protocol: parseResult[1],
		Host:     parseResult[2],
		Org:      parseResult[3],
		Repo:     parseResult[4],
		Path:     parseResult[5],
		Branch:   parseResult[6],
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
