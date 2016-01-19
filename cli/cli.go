package cli

import "github.com/tucnak/climax"

type CLI struct {
	Ctrl *Controller
}

func New(ctlr *Controller) (c *CLI) {
	c = &CLI{
		Ctrl: ctlr,
	}

	return
}

func (c *CLI) Run() {
	dcdr := climax.New("dcdr")
	dcdr.Brief = "Decider: CLI for decider feature flags."
	dcdr.Version = "stable"

	cmds := c.Commands()

	for _, cmd := range cmds {
		dcdr.AddCommand(cmd)
	}

	dcdr.Run()
}

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
					Name:     "type",
					Short:    "t",
					Usage:    `--type=[boolean|percentile]`,
					Help:     `the type of flag to set`,
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
			},

			Examples: []climax.Example{
				{
					Usecase:     `-n "flag_name" -t percentile -v 0.5 -c "the flag desc"`,
					Description: `sets a percentile flag to 50%`,
				},
				{
					Usecase:     `-n "flag_name" -t boolean -v false -c "the flag desc"`,
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
			Brief: "init the audit repo",
			Usage: `--create creates an empty audit repo and pushes to origin`,
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
	}
}
