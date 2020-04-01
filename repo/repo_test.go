package repo

import (
    "fmt"
    "io"
    "os"
    "os/exec"
    "testing"

    "github.com/stretchr/testify/assert"
)

var (
    project1 = Project { Path: "primary", Name: "filbird/v4-repo-test.git" }
    project2 = Project { Path: "dependencies/timeout", Name: "filbird/v4-timeout-test.git" }
    project3 = Project { Path: "dependencies/cron", Name: "filbird/v4-cron-test.git" }
    testManifest = Manifest { Projects: []Project { project1, project2, project3 } }
);

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestMain(m *testing.M) {
    open = func(f string) (*os.File, error) {
		return os.Open("../data/manifest.xml")
	}

    os.Exit(m.Run())
}

func TestInit(t *testing.T) {
    execCommand = fakeExecCommand
    defer func() { execCommand = exec.Command }()

    err := Init("filbird/v4-repo-test.git", "master", "")

    assert.Nil(t, err)
}

func TestSync(t *testing.T) {
    execCommand = fakeExecCommand
    defer func() { execCommand = exec.Command }()

    err := Sync()

    assert.Nil(t, err)
}

func TestParseManifestFile(t *testing.T) {
    oldOpen := open
    defer func() { open = oldOpen }()

    execCommand = fakeExecCommand
    defer func() { execCommand = exec.Command }()

    open = func(f string) (*os.File, error) {
		return os.Open("../data/manifest.xml")
	}

    manifest, err := ParseManifestFile();
    if err != nil {
        t.Errorf("ParseManifestFile() error = %q, should be nil", err)
    }

    assert.Equal(t, *manifest, testManifest)
}

func TestParseManifestFileOpenError(t *testing.T) {
    oldOpen := open
    defer func() { open = oldOpen }()

    execCommand = fakeExecCommand
    defer func() { execCommand = exec.Command }()

    open = func(f string) (*os.File, error) {
		return nil, fmt.Errorf("Spooky error")
	}

    manifest, err := ParseManifestFile();

    assert.Nil(t, manifest)
    assert.Equal(t, err, fmt.Errorf("Spooky error"))
}

func TestParseManifestFileReadAllError(t *testing.T) {
    oldReadAll := readAll
    defer func() { readAll = oldReadAll }()

    execCommand = fakeExecCommand
    defer func() { execCommand = exec.Command }()

    readAll = func(r io.Reader) ([]byte, error) {
		return []byte{}, fmt.Errorf("Spooky error")
	}

    manifest, err := ParseManifestFile();

    assert.Nil(t, manifest)
    assert.Equal(t, err, fmt.Errorf("Spooky error"))
}

func TestParseManifestUnmarshlError(t *testing.T) {
    oldUnmarshal := unmarshal
    defer func() { unmarshal = oldUnmarshal }()

    execCommand = fakeExecCommand
    defer func() { execCommand = exec.Command }()

    unmarshal = func(data []byte, v interface{}) error {
		return fmt.Errorf("Spooky error")
	}

    manifest, err := ParseManifestFile();

    assert.Nil(t, manifest)
    assert.Equal(t, err, fmt.Errorf("Spooky error"))
}

func TestFindProject(t *testing.T) {
    sourceRepo := "filbird/v4-repo-test"
    project := FindProject(&testManifest, sourceRepo);

    assert.Equal(t, *project, project1)
}

func TestFindProjectDoesNotExist(t *testing.T) {
    sourceRepo := "fake"
    project := FindProject(&testManifest, sourceRepo);

    assert.Nil(t, project)
}

// This is a fake test for mocking out exec calls.
// See https://golang.org/src/os/exec/exec_test.go and
// https://npf.io/2015/06/testing-exec-command/ for more info
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args[:]
	for i, val := range os.Args { // Should become something lke ["repo", "init"]
		args = os.Args[i:]
		if val == "--" {
			args = args[1:]
			break
		}
	}

	if len(args) >= 2 && args[0] == "repo" {
		switch args[1] {
		default:
            os.Exit(255)
        case "version":
            fmt.Println("repo version v2.4.334")
            return
		case "init":
            fmt.Println("Initializing repo manifest directory")
            return
		case "sync":
            fmt.Println("Syncing repo dependencies")
			return
		}
	}

	os.Exit(255)
}
