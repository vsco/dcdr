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
	tbl := table.New("Name", "Type", "Value", "Scope", "Updated By", "Comment").
		WithHeaderFormatter(headerFmt).
		WithFirstColumnFormatter(columnFmt)

	for _, feature := range features {
		tbl.AddRow(feature.Key, feature.FeatureType, feature.Value, feature.Scope, feature.UpdatedBy, feature.Comment)
	}

	tbl.Print()
}

func (u *UI) DrawConfig(cfg *config.Config) {
	tbl := table.New("Component", "Name", "Value", "Description").WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	tbl.AddRow("Config", "ConfigPath", config.ConfigPath(), "Path to config.hcl")
	tbl.AddRow("Defaults", "Username", cfg.Username, "The username for commits. default: `whoami`")
	tbl.AddRow("Defaults", "Namespace", cfg.Namespace, "K/V namespace")

	tbl.AddRow("Watcher", "OutputPath", cfg.Watcher.OutputPath, "File path to watch and read from")

	tbl.AddRow("Server", "Endpoint", cfg.Server.Endpoint, "The path to serve (GET '/dcdr.json')")
	tbl.AddRow("Server", "Host", cfg.Server.Host, "The server host (:8000")
	tbl.AddRow("Server", "JsonRoot", cfg.Server.JsonRoot, "JSON root node ('dcdr')")

	if cfg.GitEnabled() {
		tbl.AddRow("Git", "RepoPath", cfg.Git.RepoPath, "Audit repo location")
		tbl.AddRow("Git", "RepoURL", cfg.Git.RepoURL, "Audit repo remote origin")
	}

	if cfg.StatsEnabled() {
		tbl.AddRow("Stats", "Namespace", cfg.Stats.Namespace, "Prefix for `dcdr` change events.")
		tbl.AddRow("Stats", "Host", cfg.Stats.Host, "Statsd host ('localhost')")
		tbl.AddRow("Stats", "Port", fmt.Sprintf("%d", cfg.Stats.Port), "Statsd port (8125)")
	}

	tbl.Print()
}
