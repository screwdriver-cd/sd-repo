package main

import(
  "fmt"
  "flag"
  "os"
  "io/ioutil"

  "github.com/screwdriver-cd/sd-repo/git"
  "github.com/screwdriver-cd/sd-repo/repo"
)

// sourcePathFile is the file in which the checked out source path
// will be written to
const sourcePathFile = "sourcePath"

var (
    newGitUrl = git.New
    repoInit = repo.Init
    repoSync = repo.Sync
    parseManifestFile = repo.ParseManifestFile
    findProject = repo.FindProject
    writeFile = ioutil.WriteFile
)

// run executes the Repo workflow for 'getCheckoutCommand' in screwdriver-scm-github
// It initializes the Repo repository via `repo init`
// After parsing and validating the manifest file, it checks out dependencies via `repo sync`
// Lastly, it outputs the source repository's path to sourcePathFile
func run(manifestUrl, sourceRepo string) error {
    // Validate the manifestUrl
    gitUrl, err := newGitUrl(manifestUrl)
    if err != nil {
        return fmt.Errorf("Error validating manifest URL: %v\n", err)
    }
    checkoutUrl, branch := gitUrl.GetCloneInfo()

    // Initialize Repo repository
    err = repoInit(checkoutUrl,branch,gitUrl.Path)
    if err != nil {
        return fmt.Errorf("Error executing 'repo init': %v\n", err)
    }

    manifest, err := parseManifestFile()
    if err != nil {
        return fmt.Errorf("Error parsing manifest file: %v\n", err)
    }

    project := findProject(manifest, sourceRepo)
    if project == nil {
        return fmt.Errorf("Error: Source repo, %v, is not listed in the manifest file, %v\n", sourceRepo, gitUrl.Path)
    }

    // Checkout dependencies
    err = repoSync()
    if err != nil {
        return fmt.Errorf("Error executing 'repo sync': %v\n", err)
    }

    fmt.Printf("Writing source repository path to %v\n", sourcePathFile)
    err = writeFile(sourcePathFile, []byte(project.Path), 0444)
    if err != nil {
        return fmt.Errorf("Error writing to %v: %v\n", sourcePathFile, err)
    }

    return nil
}

func main() {
    manifestUrl := flag.String("manifestUrl", "", "URL of the repository containing the repo manifest file")
    sourceRepo := flag.String("sourceRepo", "", "Source repository")
    flag.Parse()

    err := run(*manifestUrl, *sourceRepo);
    if err != nil {
        fmt.Printf("%v", err)
        os.Exit(1)
    }
}
