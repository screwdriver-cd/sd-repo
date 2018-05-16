package repo

import (
  "encoding/xml"
  "fmt"
  "io/ioutil"
  "os"
  "os/exec"
  "strings"
)

// manifestPath is generated following successful execution of "repo init"
const manifestPath string = ".repo/manifest.xml"

var (
    execCommand = exec.Command
    open = os.Open
    readAll = ioutil.ReadAll
    unmarshal = xml.Unmarshal
)

// A Project represents a <project> XML element
type Project struct {
  Path string `xml:"path,attr"`
  Name string `xml:"name,attr"`
}

// A Manifest represents a <manifest> XML element
type Manifest struct {
  Projects []Project `xml:"project"`
}

// Init executes the "repo init" command
func Init(checkoutUrl, branch, path string) error {
    if path == "" {
        path = "default.xml"
    }

    cmd := execCommand("repo", "init", "-u", checkoutUrl, "-b", branch, "-m", path)

    fmt.Printf("repo init -u %v -b %v -m %v\n", checkoutUrl, branch, path)
    output, err := cmd.CombinedOutput()
    fmt.Println(string(output[:]))

    return err
}

// Sync executes the "repo sync" command
func Sync() error {
    cmd := execCommand("repo", "sync", "-d", "-c", "--jobs=4")

    fmt.Println("repo sync -d -c --jobs=4")
    output, err := cmd.CombinedOutput()
    fmt.Println(string(output[:]))

    return err
}

// ParseManifestFile parses the manifest file located in manifestPath
// If valid, a Manifest object is returned
func ParseManifestFile() (*Manifest, error) {
    xmlFile, err := open(manifestPath)
    if err != nil {
        return nil, err
    }
    defer xmlFile.Close()

    content, err := readAll(xmlFile)
    if err != nil {
        return nil, err
    }

    manifest := Manifest {}

    err = unmarshal([]byte(content), &manifest)
    if err != nil {
        return nil, err
    }

    return &manifest, err
}

// FindProject returns a Project in Manifest whose Name matches sourceRepo
// If no project is found, nil is returned
func FindProject(manifest *Manifest, sourceRepo string) *Project {
    for _, project := range manifest.Projects {
        projectName := strings.TrimSuffix(project.Name, ".git")
        sourceRepo = strings.TrimSuffix(sourceRepo, ".git")
        if projectName == sourceRepo {
            return &project
        }
    }
    return nil
}
