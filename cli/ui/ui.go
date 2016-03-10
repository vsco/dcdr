package ui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/vsco/dcdr/cli/models"
	"github.com/vsco/dcdr/config"
)

var (
	headerFmt = color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt = color.New(color.FgYellow).SprintfFunc()
)

type UI struct{}

func New() (u *UI) {
	u = &UI{}

	return
}

func (u *UI) DrawFeatures(features models.Features) {
	color.NoColor = false
	tbl := table.New("Name", "Type", "Value", "Comment", "Scope", "Updated By").
		WithHeaderFormatter(headerFmt).
		WithFirstColumnFormatter(columnFmt)

	for _, feature := range features {
		tbl.AddRow(feature.Key, feature.FeatureType, feature.Value, feature.Comment, feature.Scope, feature.UpdatedBy)
	}

	tbl.Print()
}

func (u *UI) DrawConfig(cfg *config.Config) {
	tbl := table.New("Component", "Name", "Value", "Description").WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	tbl.AddRow("Defaults", "Username", cfg.Username, "The username for audit commits. `whoami` unless set in `config.hcl`.")
	tbl.AddRow("Defaults", "Namespace", cfg.Namespace, "Consul KV namespace to use for feature flags")

	tbl.AddRow("Watcher", "OutputPath", cfg.Watcher.OutputPath, "Path to the file written by watch and read by `Client`.")

	tbl.AddRow("Server", "Endpoint", cfg.Server.Endpoint, "The path at which to serve feature flags. ('/dcdr.json')")
	tbl.AddRow("Server", "Host", cfg.Server.Host, "The host used by the server. (:8000")
	tbl.AddRow("Server", "JsonRoot", cfg.Server.JsonRoot, "Root json node served by `Endpoint`. ('dcdr')")

	if cfg.GitEnabled() {
		tbl.AddRow("Git", "RepoPath", cfg.Git.RepoPath, "Location on disk for the audit repo.")
		tbl.AddRow("Git", "RepoURL", cfg.Git.RepoURL, "Remote origin for the autdit repo.")
	}

	if cfg.StatsEnabled() {
		tbl.AddRow("Stats", "Namespace", cfg.Stats.Namespace, "Namespace prefix for `dcdr` change events.")
		tbl.AddRow("Stats", "Host", cfg.Stats.Host, "Statsd host ('localhost')")
		tbl.AddRow("Stats", "Port", fmt.Sprintf("%s", cfg.Stats.Port), "Statsd port (8125)")
	}

	tbl.Print()
}
