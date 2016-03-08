package cli

import "github.com/tucnak/climax"

const Version = "0.1"

// CLI main CLI runner
type CLI struct {
	Ctrl *Controller
}

// New init a new CLI
func New(ctlr *Controller) (c *CLI) {
	c = &CLI{
		Ctrl: ctlr,
	}

	return
}

// Run bind command and run
func (c *CLI) Run() {
	dcdr := climax.New("dcdr")
	dcdr.Brief = "Decider: CLI for decider feature flags."
	dcdr.Version = Version

	cmds := c.Commands()

	for _, cmd := range cmds {
		dcdr.AddCommand(cmd)
	}

	dcdr.Run()
}

// Commands slice of all commands
func (c *CLI) Commands() []climax.Command {
	return []climax.Command{
		{
			Name:  "list",
			Brief: "list all feature flags",
			Usage: `[-p=] "<prefix>" list all flags with a matching prefix`,
			Help:  `Lists all feature flags. Use -p to match flags by a prefix.`,

			Flags: []climax.Flag{
				{
					Name:     "prefix",
					Short:    "p",
					Usage:    `--prefix="<flag_name>"`,
					Help:     `List only flags with matching prefix.`,
					Variable: true,
				},
				{
					Name:     "scope",
					Short:    "s",
					Usage:    `--scope="<flag_scope>"`,
					Help:     `List only flags within a scope.`,
					Variable: true,
				},
			},

			Examples: []climax.Example{
				{
					Usecase:     `-p "flag_"`,
					Description: `Matches 'flag_name'`,
				},
			},

			Handle: c.Ctrl.List,
		},
		{
			Name:  "set",
			Brief: "create or update a feature flag",
			Usage: `set -name flag_name -type [boolean|percentile] -value [0.0-1.0|true/false] -comment "flag description"`,
			Help:  `set creates or updates a feature flag.`,

			Flags: []climax.Flag{
				{
					Name:     "name",
					Short:    "n",
					Usage:    `--name="flag_name"`,
					Help:     `the name of the falg to set`,
					Variable: true,
				},
				{
					Name:     "value",
					Short:    "v",
					Usage:    `--value=0.0-1.0 or true|false`,
					Help:     `the value of the flag`,
					Variable: true,
				},
				{
					Name:     "comment",
					Short:    "c",
					Usage:    `--comment="flag description"`,
					Help:     `an optional comment or description`,
					Variable: true,
				},
				{
					Name:     "scope",
					Short:    "s",
					Usage:    `--scope="flag scope"`,
					Help:     `an optional scope to nest the flag within`,
					Variable: true,
				},
			},

			Examples: []climax.Example{
				{
					Usecase:     `-n "flag_name" -v 0.5 -c "the flag desc"`,
					Description: `sets a percentile flag to 50%`,
				},
				{
					Usecase:     `-n "flag_name" -v false -c "the flag desc"`,
					Description: `sets a boolean flag to false`,
				},
			},

			Handle: c.Ctrl.Set,
		},
		{
			Name:  "delete",
			Brief: "delete a feature flag",
			Usage: `[-n=] "<name>" delete flag with matching name`,
			Help:  `Delete a feature flag matching --name`,

			Flags: []climax.Flag{
				{
					Name:     "name",
					Short:    "n",
					Usage:    `--name="<flag_name>"`,
					Help:     `Name of the flag to delete`,
					Variable: true,
				},
				{
					Name:     "scope",
					Short:    "s",
					Usage:    `--scope="flag scope"`,
					Help:     `an optional scope to delete the flag from.`,
					Variable: true,
				},
			},

			Examples: []climax.Example{
				{
					Usecase:     `-n "flag_name"`,
					Description: `Deletes 'flag_name'`,
				},
			},

			Handle: c.Ctrl.Delete,
		},
		{
			Name:  "init",
			Brief: "init the audit repository",
			Usage: `--create creates an empty audit repository and pushes to origin`,
			Help: `Clones the RepoUrl into the RepoPath from ~/.dcdr. Creates a new
			repo if --create is passed.`,

			Flags: []climax.Flag{
				{
					Name:     "create",
					Short:    "c",
					Usage:    `--create`,
					Help:     `Create a new empty repo`,
					Variable: false,
				},
			},

			Examples: []climax.Example{
				{
					Usecase:     ``,
					Description: `clone RepoUrl into the RepoPath from ~/.dcdr`,
				},
				{
					Usecase:     `--create`,
					Description: `create a new empty repo`,
				},
			},

			Handle: c.Ctrl.Init,
		},
		{
			Name:  "info",
			Brief: "display config info",
			Usage: ``,
			Help:  `Displays current settings from config.hcl`,

			Handle: c.Ctrl.Info,
		},
		{
			Name:  "server",
			Brief: "start dcdr http server",
			Usage: ``,
			Help:  `Starts the dcdr http server`,

			Handle: c.Ctrl.Serve,
		},
		{
			Name:  "import",
			Brief: "import json from STDIN",
			Usage: ``,
			Help:  `Imports JSON feature flags`,

			Handle: c.Ctrl.Import,
		},
		{
			Name:  "watch",
			Brief: "watch the dcdr namespace for changes",
			Usage: ``,
			Help:  `watch the dcdr namespace for changes and write the JSON file used by clients`,

			Handle: c.Ctrl.Watch,
		},
	}
}
