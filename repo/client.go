package repo

import (
	"os"
	"os/exec"

	"fmt"
	"strings"

	"encoding/json"
	"io/ioutil"

	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
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
	Config *config.Config
}

func New(c *config.Config) (g *Git) {
	g = &Git{
		Config: c,
	}

	return
}

func (g *Git) Create() error {
	err := os.MkdirAll(g.Config.Git.RepoPath, DefaultPerms)

	if err != nil {
		return fmt.Errorf("failed to create repo: %v\n", err)
	}

	fp := fmt.Sprintf("%s/%s", g.Config.Git.RepoPath, DefaultConfigFileName)
	err = ioutil.WriteFile(fp, []byte{}, DefaultPerms)

	if err != nil {
		return fmt.Errorf("failed to create %s: %v\n", DefaultConfigFileName, err)
	}

	cmd := exec.Command(GitExec(), "init")
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to init repo into %s\n", g.Config.Git.RepoPath)
	}

	cmd = exec.Command(GitExec(), "add", ".")
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to add %s into %s\n", fp, g.Config.Git.RepoPath)
	}

	msg := "Initializing decider repo"

	cmd = exec.Command(GitExec(), "commit", "-am", msg)
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to commit into %s\n", g.Config.Git.RepoPath)
	}

	cmd = exec.Command(GitExec(), "remote", "add", "origin", g.Config.Git.RepoURL)
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to add origin %s\n", g.Config.Git.RepoURL)
	}

	cmd = exec.Command(GitExec(), "push", "origin", "master")
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to push to %s\n", g.Config.Git.RepoURL)
	}

	fmt.Printf("created %s and pushed to %s sucessfully\n", g.Config.Git.RepoPath, g.Config.Git.RepoURL)

	return nil
}

func (g *Git) Clone() error {
	_, err := exec.Command(GitExec(), "clone", g.Config.Git.RepoURL, g.Config.Git.RepoPath).Output()

	if err != nil {
		return fmt.Errorf("could not checkout %s into %s\n", g.Config.Git.RepoURL, g.Config.Git.RepoPath)
	}

	fmt.Printf("cloned %s into %s\n", g.Config.Git.RepoURL, g.Config.Git.RepoPath)

	return nil
}

func (g *Git) nothingToCommit(msg []byte) bool {
	return strings.Contains(string(msg[:]), "nothing to commit")
}

func (g *Git) Pull() error {
	cmd := exec.Command(GitExec(), "pull", "origin", "master")
	cmd.Dir = g.Config.Git.RepoPath
	bts, err := cmd.Output()

	if err != nil {
		return fmt.Errorf(string(bts[:]))
	}

	return nil
}

func (g *Git) Commit(features models.Features, msg string) error {
	if !g.Config.UseGit() {
		return nil
	}

	if err := g.Pull(); err != nil {
		return fmt.Errorf("could not pull from %s", g.Config.Git.RepoURL)
	}

	bts, _ := json.MarshalIndent(features, "", "  ")

	fp := fmt.Sprintf("%s/%s", g.Config.Git.RepoPath, DefaultConfigFileName)
	err := ioutil.WriteFile(fp, bts, DefaultPerms)

	if err != nil {
		return fmt.Errorf("could not write change to %s\n", fp)
	}

	msg = fmt.Sprintf("%s %s", g.Config.Username, msg)

	cmd := exec.Command(GitExec(), "commit", "-am", msg)
	cmd.Dir = g.Config.Git.RepoPath
	out, err := cmd.Output()

	if g.nothingToCommit(out) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("could not commit change to %s %s\n", g.Config.Git.RepoPath, string(out[:]))
	}

	cmd = exec.Command(GitExec(), "push", "origin", "master")
	cmd.Dir = g.Config.Git.RepoPath
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to push to %s\n", g.Config.Git.RepoURL)
	}

	return nil
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
