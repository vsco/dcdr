package cli

import "github.com/tucnak/climax"

const Version = "0.2"

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

// Run bind commands and run
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
			Help: `


	Lists feature flags. Use --prefix to match flags by a prefix and --scope to match
	only flags within a given scope.`,

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
			Help: `


	Set creates or updates a feature flag. Flags are set in the <Namespace>/default
	if no --scope param is provided. Scopes can have an arbitrary depth and result in
	a nesting of the feature map according to their slashes.

	Example:

	# set a key to 100% for only the US country code
	dcdr set -n new-signup-flow -v 1.0 -s country-codes/us

	Would be reflected when <FeatureMapPath> is written by 'dcdr watch' as the following.

	{
		"dcdr": {
			"info": {
				"current_sha": "faf9b666c0a51e66bc36828f819f0497720d4215"
			},
			"features": {
				"default": {
					"new-signup-flow": 0
				},
				"country-codes": {
					"us": {
						"new-signup-flow": 1.0
					}
				}
			}
		}
	}

	Clients can then read this value by using WithScopes. If no value is found within
	the provided scope, the client will fall back to the 'default' namespace value
	or return false when the key cannot be found.

	d := dcdr.NewDefault().WithScopes("cc/us")

	fmt.Printf("%t", d.IsAvailableForId("new-signup-flow", <unint64>))
	=> true

	If the audit repo has been configured in config.hcl, dcdr
	will export the full feature set and write it to <Git:RepoPath> and then
	attempt to commit and push the changeset to <Git:RepoURL>. If the commit is successful
	the 'git rev-parse HEAD' will be set into '<Namespace/info/current_sha>'.
	`,

			Flags: []climax.Flag{
				{
					Name:     "name",
					Short:    "n",
					Usage:    `--name="flag_name"`,
					Help:     `the name of the flag to set`,
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
			Usage: `[-n or --name]= "<name>" delete flag with matching name`,
			Help: `


	Delete a feature flag matching --name. Use --scope to delete a flag within
	a given scope. If the audit repo has been configured in config.hcl, dcdr
	will export the full feature set and write it to <Git:RepoPath> and then
	attempt to commit and push the changeset to <Git:RepoURL>.`,

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
			Usage: `--create creates default config.hcl and empty audit repository and pushes to origin`,
			Help: `


	By default Decider looks in /etc/dcdr for config.hcl. This can be overridden by setting the DCDR_CONFIG_DIR
	environment variable. If no config.hcl file is found in this directory, init will attempt to create one. This file
	contains an example config with all settings commented out. If no /etc/dcdr directory exists you will need to
	create this yourself. Depending on your permissions, try the following.

	sudo mkdir /etc/dcdr
	sudo chown $(whoami) /etc/dcdr

	If a repository has been configured init clones the <Git:RepoUrl> into the <Git:RepoPath> from config.hcl.
	To create a new repo with an empty decider.json pass the --create flag.`,

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
			Help: `

	Starts the dcdr http server on <Server:Host>/<Sever:Endpoint>.
	The default is 'http://localhost:8000/dcdr.json'. The root json node can also be
	set by changing the <Server:JsonRoot> value.

	Example:

	{
		"<Server:JsonRoot>": {
			"some-feature":true,
			"some-other-feature": 0.5
		}
	}
	`,

			Handle: c.Ctrl.Serve,
		},
		{
			Name:  "import",
			Brief: "import json from STDIN",
			Usage: ``,
			Help: `

	Imports JSON feature flags from a flat JSON KV structure such as...

	{
		"some-feature":true,
		"some-other-feature": 0.5
	}

	These KV pairs are set into the <Namespace>/default scope unless a --scope
	param is provided.`,

			Flags: []climax.Flag{
				{
					Name:     "scope",
					Short:    "s",
					Usage:    `--scope="some-scope"`,
					Help:     `scope to import the KVs into`,
					Variable: true,
				},
			},

			Handle: c.Ctrl.Import,
		},
		{
			Name:  "watch",
			Brief: "watch the dcdr namespace for changes",
			Usage: ``,
			Help: `


	Watches the consul KV store '<Namespace>' for changes. When changes are
	observed, dcdr writes the entire keyspace to the <FeatureMapPath> from
	config.hcl as a nested JSON hash. These events trigger a FeatureMap
	update within any dcdr.Client observing that file.

	Example Output:

	{
		"dcdr": {
			"info": {
				"current_sha": "faf9b666c0a51e66bc36828f819f0497720d4215"
			},
			"features": {
				"default": {
					"new-signup-flow": 0
				},
				"country-codes": {
					"us": {
						"new-signup-flow": 1.0
					}
				}
			}
		}
	}`,

			Handle: c.Ctrl.Watch,
		},
	}
}
