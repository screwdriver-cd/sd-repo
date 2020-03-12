package main

import(
    "fmt"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/screwdriver-cd/sd-repo/git"
    "github.com/screwdriver-cd/sd-repo/repo"
)

const (
    testManifestUrl string= "git@github.com:filbird/v4-repo-test.git/default.xml"
    testSourceRepo string= "filbird/v4-repo-test.git"
)

var (
    testProject = repo.Project { Path: "primary", Name: "filbird/v4-repo-test.git" }
    testManifest = repo.Manifest { Projects: []repo.Project { testProject } }
    testGitUrl = git.GitUrl {
        Host:   "github.com",
        Org:    "filbird",
        Repo:   "v4-repo-test",
        Path:   "default.xml",
        Branch: "master",
    }
);

func TestMain(m *testing.M) {
    // Mcok out functions
    newGitUrl  = func(testGitUrlStr string) (*git.GitUrl, error) {return &testGitUrl, nil }
    repoInit = func(checkoutUrl, branch, path string) error { return nil }
    parseManifestFile = func() (*repo.Manifest, error) { return &testManifest, nil }
    repoSync = func() error { return nil }
    parseManifestFile = func() (*repo.Manifest, error) { return &testManifest, nil }
    findProject = func(testManifest *repo.Manifest, testSourceRepo string) *repo.Project { return &testProject }
    writeFile = func(string, []byte, os.FileMode) error { return nil }

    os.Exit(m.Run())
}

func TestRun(t *testing.T) {
    err := run(testManifestUrl, testSourceRepo)
    assert.Nil(t, err)
}

func TestRunNewGitUrlError(t *testing.T) {
    oldNewGitUrl := newGitUrl
    defer func() { newGitUrl = oldNewGitUrl }()
    newGitUrl  = func(testGitUrlStr string) (*git.GitUrl, error) {
        return nil, fmt.Errorf("Spooky error")
    }

    err := run(testManifestUrl, testSourceRepo)

    expectedError := fmt.Errorf("Error validating manifest URL: Spooky error\n")
    assert.Equal(t, err, expectedError)
}

func TestRunRepoInitError(t *testing.T) {
    oldRepoInit := repoInit
    defer func() { repoInit = oldRepoInit }()
    repoInit = func(checkoutUrl, branch, path string) error {
        return fmt.Errorf("Spooky error")
    }

    err := run(testManifestUrl, testSourceRepo)

    expectedError := fmt.Errorf("Error executing 'repo init': Spooky error\n")
    assert.Equal(t, err, expectedError)
}

func TestRunParseManifestFileError(t *testing.T) {
    oldParseManifestFile := parseManifestFile
    defer func() { parseManifestFile = oldParseManifestFile }()
    parseManifestFile = func() (*repo.Manifest, error) {
        return &testManifest, fmt.Errorf("Spooky error")
    }

    err := run(testManifestUrl, testSourceRepo)

    expectedError := fmt.Errorf("Error parsing manifest file: Spooky error\n")
    assert.Equal(t, err, expectedError)
}

func TestRunFindProjectError(t *testing.T) {
    oldFindProject := findProject
    defer func() { findProject = oldFindProject }()
    findProject = func(testManifest *repo.Manifest, testSourceRepo string) *repo.Project {
        return nil
    }

    err := run(testManifestUrl, testSourceRepo)

    expectedError := fmt.Errorf("Error: Source repo, %v, is not listed in the manifest file, %v\n", testSourceRepo, testGitUrl.Path)
    assert.Equal(t, err, expectedError)
}

func TestRunRepoSyncError(t *testing.T) {
    oldRepoSync := repoSync
    defer func() { repoSync = oldRepoSync }()
    repoSync = func() error { return fmt.Errorf("Spooky error") }

    err := run(testManifestUrl, testSourceRepo)

    expectedError := fmt.Errorf("Error executing 'repo sync': Spooky error\n")
    assert.Equal(t, err, expectedError)
}

func TestRunWriteFileError(t *testing.T) {
    oldWriteFile := writeFile
    defer func() { writeFile = oldWriteFile }()
    writeFile = func(string, []byte, os.FileMode) error {
        return fmt.Errorf("Spooky error")
    }

    err := run(testManifestUrl, testSourceRepo)

    expectedError := fmt.Errorf("Error writing to %v: Spooky error\n", sourcePathFile)
    assert.Equal(t, err, expectedError)
}
