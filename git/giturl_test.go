package git

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseGitUrlSuccess(t *testing.T) {
	var gitUrl, _ = New("git@github.com:screwdriver-cd/sd-repo.git/model/giturl_test.go#test")
	assert.Equal(t, "github.com", gitUrl.Host)
	assert.Equal(t, "screwdriver-cd", gitUrl.Org)
	assert.Equal(t, "sd-repo", gitUrl.Repo)
	assert.Equal(t, "model/giturl_test.go", gitUrl.Path)
	assert.Equal(t, "test", gitUrl.Branch)

	var gitUrlNoPathNoBranch, _ = New("git@github.com:screwdriver-cd/sd-repo.git")
	assert.Equal(t, "github.com", gitUrlNoPathNoBranch.Host)
	assert.Equal(t, "screwdriver-cd", gitUrlNoPathNoBranch.Org)
	assert.Equal(t, "sd-repo", gitUrlNoPathNoBranch.Repo)
	assert.Equal(t, "", gitUrlNoPathNoBranch.Path)
	assert.Equal(t, "master", gitUrlNoPathNoBranch.Branch)

	var gitUrlWeirdBranch, _ = New("git@gitgit.com:screwdriver-cd/sd-repo.git##test")
	assert.Equal(t, "gitgit.com", gitUrlWeirdBranch.Host)
	assert.Equal(t, "screwdriver-cd", gitUrlWeirdBranch.Org)
	assert.Equal(t, "sd-repo", gitUrlWeirdBranch.Repo)
	assert.Equal(t, "", gitUrlWeirdBranch.Path)
	assert.Equal(t, "#test", gitUrlWeirdBranch.Branch)

	var gitUrlWeirdPath, _ = New("git@github.com:screwdriver-cd//sd-repo.git//a/bb/c.xml")
	assert.Equal(t, "github.com", gitUrlWeirdPath.Host)
	assert.Equal(t, "screwdriver-cd", gitUrlWeirdPath.Org)
	assert.Equal(t, "sd-repo", gitUrlWeirdPath.Repo)
	assert.Equal(t, "a/bb/c.xml", gitUrlWeirdPath.Path)
	assert.Equal(t, "master", gitUrlWeirdPath.Branch)
}

func TestParseGitUrlError(t *testing.T) {
	var gitUrlBad1, err1 = New("git@github.com::screwdriver-cd/sd-repo.git/model/giturl_test.go#test")
	assert.Nil(t, gitUrlBad1)
	if assert.Error(t, err1, "should return error on invalid git config") {
		assert.Equal(t, "Error: not a valid git url git@github.com::screwdriver-cd/sd-repo.git/model/giturl_test.go#test", err1.Error())
	}

	var gitUrlBad2, err2 = New("git@github.com:sd-repo.git/model/giturl_test.go#test")
	assert.Nil(t, gitUrlBad2)
	if assert.Error(t, err2, "should return error on invalid git config") {
		assert.Equal(t, "Error: not a valid git url git@github.com:sd-repo.git/model/giturl_test.go#test", err2.Error())
	}

	var gitUrlBad3, err3 = New("git@github.com:a/b/model/giturl_test.git")
	assert.Nil(t, gitUrlBad3)
	if assert.Error(t, err3, "should return error on invalid git config") {
		assert.Equal(t, "Error: not a valid git url git@github.com:a/b/model/giturl_test.git", err3.Error())
	}

	var gitUrlBad4, err4 = New("github.com:a/b.git#branch")
	assert.Nil(t, gitUrlBad4)
	if assert.Error(t, err4, "should return error on invalid git config") {
		assert.Equal(t, "Error: not a valid git url github.com:a/b.git#branch", err4.Error())
	}
}

func TestGetCloneInfo(t *testing.T) {
	var gitUrl1, err1 = New("git@gitgit.com:screwdriver-cd/sd-repo.git##test")
	assert.Nil(t, err1)
	var gitCloneUrl1, branch1 = gitUrl1.GetCloneInfo()
	assert.Equal(t, "git@gitgit.com:screwdriver-cd/sd-repo.git", gitCloneUrl1)
	assert.Equal(t, "#test", branch1)

	var gitUrl2, err2 = New("git@gitgit.com:screwdriver-cd/sd-repo2.git/blah/blhablha")
	assert.Nil(t, err2)
	var gitCloneUrl2, branch2 = gitUrl2.GetCloneInfo()
	assert.Equal(t, "git@gitgit.com:screwdriver-cd/sd-repo2.git", gitCloneUrl2)
	assert.Equal(t, "master", branch2)
}
