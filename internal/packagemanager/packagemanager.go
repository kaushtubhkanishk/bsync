package packagemanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Package struct {
	BinPath        string   `yaml:"bin_path"`
	BuildBinPath   string   `yaml:"build_bin_path"`
	BuildSteps     []string `yaml:"build_steps"`
	CurrentVersion string   `yaml:"current_version"`
	GitRepo        string   `yaml:"git_repo"`
	GitAuthor      string   `yaml:"git_author"`
	Name           string   `yaml:"name"`
	Path           string   `yaml:"path"`
	LatestVersion  string   `yaml:"-"`
}

const gitRestUrl = "https://api.github.com/repos"

func (p *Package) FetchLatestVersion() error {
	latestVersionUrl, err := url.JoinPath(gitRestUrl, p.GitAuthor, p.GitRepo, "releases", "latest")
	if err != nil {
		errMsg := fmt.Sprintf("error formulating git remote url: %v", err.Error())
		return errors.New(errMsg)
	}

	httpClient := http.Client{
		Timeout: 120 * time.Second,
	}

	resp, err := httpClient.Get(latestVersionUrl)
	if err != nil {
		errMsg := fmt.Sprintf("error connecting to remote git server: %v", err.Error())
		return errors.New(errMsg)
	}

	var respData map[string]any

	err = json.NewDecoder(resp.Body).Decode(&respData)
	defer resp.Body.Close()
	if err != nil {
		errMsg := fmt.Sprintf("error reading git server response: %v", err.Error())
		return errors.New(errMsg)
	}

	latestVersion, ok := respData["tag_name"]
	if !ok {
		errMsg := fmt.Sprintf("tag_name doesn't exist in the response")
		return errors.New(errMsg)
	}

	p.LatestVersion, ok = latestVersion.(string)
	if !ok {
		errMsg := fmt.Sprintf("tag_name is not of type string")
		return errors.New(errMsg)
	}

	return nil
}

func (p *Package) Update() error {
	err := p.checkoutLatestTag()
	if err != nil {
		errMsg := fmt.Sprintf("error checking out latest git tag: %v", err.Error())
		return errors.New(errMsg)
	}

	err = p.executeBuildSteps()
	if err != nil {
		errMsg := fmt.Sprintf("error running build: %v", err.Error())
		return errors.New(errMsg)
	}

	err = p.copyToBinPath()
	if err != nil {
		errMsg := fmt.Sprintf("error copying built binary to bin path: %v", err.Error())
		return errors.New(errMsg)
	}

	p.CurrentVersion = p.LatestVersion

	return nil
}

func (p *Package) checkoutLatestTag() error {
	outFetch, err := exec.Command("git", "-C", p.Path, "fetch", "--tags").Output()
	if err != nil {
		errMsg := fmt.Sprintf("error running git fetch: %v", err.Error())
		return errors.New(errMsg)
	}
	fmt.Println(string(outFetch))

	outCheckout, err := exec.Command("git", "-C", p.Path, "checkout", p.LatestVersion).Output()
	if err != nil {
		errMsg := fmt.Sprintf("error running git checkout: %v", err.Error())
		return errors.New(errMsg)
	}
	fmt.Println(string(outCheckout))

	fmt.Println("Successfully checked out new version: ", p.LatestVersion)

	return nil
}

func (p *Package) executeBuildSteps() error {
	fmt.Println("Running build")
	for i, step := range p.BuildSteps {
		p.BuildSteps[i] = strings.ReplaceAll(step, "{path}", p.Path)
		fmt.Println(p.BuildSteps[i])
	}

	for _, step := range p.BuildSteps {
		fmt.Printf("Step: %v\n", step)
		fields := strings.Fields(step)

		out, err := exec.Command(fields[0], fields[1:]...).Output()
		if err != nil {
			errMsg := fmt.Sprintf("error running build step: %v, error: %v", step, err.Error())
			return errors.New(errMsg)
		}
		fmt.Println(string(out))
	}
	return nil
}

func (p *Package) copyToBinPath() error {
	sourceFile, err := os.OpenFile(p.BuildBinPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		errMsg := fmt.Sprintf("error opening the built binary: %v", err.Error())
		return errors.New(errMsg)
	}
	defer sourceFile.Close()

	destinationFile, err := os.OpenFile(p.BinPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		errMsg := fmt.Sprintf("error opening the binary: %v", err.Error())
		return errors.New(errMsg)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		errMsg := fmt.Sprintf("error copying the file to bin path: %v", err.Error())
		return errors.New(errMsg)
	}

	err = destinationFile.Sync()
	if err != nil {
		errMsg := fmt.Sprintf("failed to sync binary to disc: %v", err.Error())
		return errors.New(errMsg)
	}

	sourceInfo, err := os.Stat(p.BuildBinPath)
	if err != nil {
		errMsg := fmt.Sprintf("failed to stat built binary: %v", err.Error())
		return errors.New(errMsg)
	}

	err = os.Chmod(p.BinPath, sourceInfo.Mode())
	if err != nil {
		errMsg := fmt.Sprintf("failed to change permissions of the binary: %v", err.Error())
		return errors.New(errMsg)
	}

	fmt.Println("Successfully copied built binary to bin path")
	fmt.Println(p.BinPath)

	return nil
}
