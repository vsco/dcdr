package git

import (
	"os"
	"os/exec"

	"fmt"
	"strings"

	"encoding/json"
	"io/ioutil"

	"github.com/vsco/decider-cli/models"
)

const (
	DefaultConfigFileName = "decider.json"
	DefaultPerms          = 0755
)

var (
	gitExec = ""
)

func GitExec() string {
	if gitExec != "" {
		return gitExec
	}

	output, err := exec.Command("which", "git").Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gitExec = strings.TrimSpace(string(output[:]))

	return gitExec
}

type Git struct {
	Config *models.Config
}

func New(c *models.Config) (g *Git) {
	g = &Git{
		Config: c,
	}

	return
}

func (g *Git) Create() {
	err := os.MkdirAll(g.Config.Git.RepoPath, DefaultPerms)

	if err != nil {
		fmt.Printf("failed to create repo: %v\n", err)
		os.Exit(1)
	}

	fp := fmt.Sprintf("%s/%s", g.Config.Git.RepoPath, DefaultConfigFileName)
	err = ioutil.WriteFile(fp, []byte{}, DefaultPerms)

	if err != nil {
		fmt.Printf("failed to create %s: %v\n", DefaultConfigFileName, err)
		os.Exit(1)
	}

	cmd := exec.Command(GitExec(), "init")
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		fmt.Printf("failed to init repo into %s\n", g.Config.Git.RepoPath)
		os.Exit(1)
	}

	cmd = exec.Command(GitExec(), "add", ".")
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		fmt.Printf("failed to add %s into %s\n", fp, g.Config.Git.RepoPath)
		os.Exit(1)
	}

	msg := "Initializing decider repo"

	cmd = exec.Command(GitExec(), "commit", "-am", msg)
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		fmt.Printf("failed to commit into %s\n", g.Config.Git.RepoPath)
		os.Exit(1)
	}

	cmd = exec.Command(GitExec(), "remote", "add", "origin", g.Config.Git.RepoURL)
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		fmt.Printf("failed to add origin %s\n", g.Config.Git.RepoURL)
		os.Exit(1)
	}

	cmd = exec.Command(GitExec(), "push", "origin", "master")
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		fmt.Printf("failed to push to %s\n", g.Config.Git.RepoURL)
		os.Exit(1)
	}

	fmt.Printf("created %s and pushed to %s sucessfully\n", g.Config.Git.RepoPath, g.Config.Git.RepoURL)
}

func (g *Git) Clone() {
	_, err := exec.Command(GitExec(), "clone", g.Config.Git.RepoURL, g.Config.Git.RepoPath).Output()

	if err != nil {
		fmt.Printf("could not checkout %s into %s\n", g.Config.Git.RepoURL, g.Config.Git.RepoPath)
		os.Exit(1)
	}

	fmt.Printf("cloned %s into %s\n", g.Config.Git.RepoURL, g.Config.Git.RepoPath)
}

func (g *Git) nothingToCommit(msg []byte) bool {
	return strings.Contains(string(msg[:]), "nothing to commit")
}

func (g *Git) Commit(features models.Features, msg string) {
	bts, _ := json.MarshalIndent(features, "", "  ")

	fp := fmt.Sprintf("%s/%s", g.Config.Git.RepoPath, DefaultConfigFileName)
	err := ioutil.WriteFile(fp, bts, DefaultPerms)

	if err != nil {
		fmt.Printf("could not write change to %s\n", fp)
		os.Exit(1)
	}

	msg = fmt.Sprintf("%s %s", g.Config.Username, msg)

	cmd := exec.Command(GitExec(), "commit", "-am", msg)
	cmd.Dir = g.Config.Git.RepoPath
	out, err := cmd.Output()

	if g.nothingToCommit(out) {
		return
	}

	if err != nil {
		fmt.Printf("could not commit change to %s %+v\n", g.Config.Git.RepoPath, string(out[:]))
		os.Exit(1)
	}

	cmd = exec.Command(GitExec(), "push", "origin", "master")
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		fmt.Printf("failed to push to %s\n", g.Config.Git.RepoURL)
		os.Exit(1)
	}
}

func (g *Git) RepoExists() bool {
	_, err := os.Stat(g.Config.Git.RepoPath + "/.git")

	if err != nil {
		return false
	}

	return true
}

func (g *Git) Init() {
	if g.Config.UseGit() && !g.RepoExists() {
		g.Clone()
	}
}
