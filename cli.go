package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"

	"flag"

	"strconv"

	"io/ioutil"

	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/hcl"
	"github.com/vsco/decider-cli/client"
	"github.com/vsco/decider-cli/git"
	"github.com/vsco/decider-cli/models"
	"github.com/vsco/decider-cli/ui"
)

const Usage = `
list	list all keys
set	set a key
delete	remove a key
`

type CLI struct {
	action string
	client *client.Client
	repo   *git.Git
}

func NewCLI(client *client.Client, g *git.Git) (c *CLI) {
	c = &CLI{
		action: os.Args[1],
		client: client,
		repo:   g,
	}

	return
}

func (c *CLI) Run() {
	switch c.action {
	case "list":
		list := flag.NewFlagSet("list", flag.ExitOnError)
		prefix := list.String("prefix", "", "search prefix")

		list.Parse(os.Args[2:])

		features, err := c.client.List(*prefix)

		if len(features) == 0 {
			fmt.Printf("No features found in namespace: %s.\n", c.client.Namespace)
			os.Exit(0)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		ui.New().DrawTable(features)

	case "set":
		set := flag.NewFlagSet("set", flag.ExitOnError)
		name := set.String("name", "", "the feature name")
		ft := set.String("type", "", "the feature type [percentile,boolean]")
		val := set.String("value", "", "the feature value")
		cmt := set.String("comment", "", "additional comment")

		set.Parse(os.Args[2:])

		msg := fmt.Sprintf("set %s to %s", *name, *val)

		if *name == "" {
			set.PrintDefaults()
			os.Exit(0)
		}

		var ftc models.FeatureType

		existing, _ := c.client.Get(*name)

		if existing != nil {
			ftc = existing.FeatureType
		} else {
			ftc = models.GetFeatureType(*ft)
		}

		switch ftc {
		case models.Percentile:
			f, err := strconv.ParseFloat(*val, 64)

			if err != nil {
				fmt.Println("invalid -value format. use -value=[0.0-1.0]")
				os.Exit(2)
			}

			c.client.SetPercentile(*name, f, *cmt)
		case models.Boolean:
			f, err := strconv.ParseBool(*val)

			if err != nil {
				fmt.Println("invalid -value format. use -value=[true,false]")
				os.Exit(2)
			}

			c.client.SetBoolean(*name, f, *cmt)
		default:
			fmt.Printf("%q is not valid type.\n", *ft)
			os.Exit(2)
		}

		features, err := c.client.List("")

		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		c.repo.Commit(features, msg)

		fmt.Printf("set %s to %s.\n", *name, *val)
		os.Exit(0)

	case "delete":
		set := flag.NewFlagSet("delete", flag.ExitOnError)
		n := set.String("name", "", "the feature name")

		set.Parse(os.Args[2:])

		if *n == "" {
			set.PrintDefaults()
			os.Exit(0)
		}

		err := c.client.Delete(*n)

		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		features, err := c.client.List("")

		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		msg := fmt.Sprintf("deleted %s", *n)
		c.repo.Commit(features, msg)

		fmt.Printf("Deleted feature '%s'\n", *n)
	case "init":
		if c.repo.Config.UseGit() {
			yn := prompt("create decider audit repository now? y/n")

			if yn == "y" {
				c.repo.Create()
			}
		}
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}

func configPath() string {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir + "/.dcdr"
}

func prompt(q string) string {
	fmt.Println(q)
	reader := bufio.NewReader(os.Stdin)
	resp, _ := reader.ReadString('\n')

	return strings.TrimSpace(resp)
}

func readConfig() *models.Config {
	bts, err := ioutil.ReadFile(configPath())

	if err != nil {
		fmt.Printf("Could not read %s", configPath())
		os.Exit(1)
	}

	var cfg *models.Config

	err = hcl.Decode(&cfg, string(bts[:]))

	if err != nil {
		fmt.Printf("parse error %+v", err)
		os.Exit(1)
	}

	return cfg
}

func loadConfig() *models.Config {
	if _, err := os.Stat(configPath()); err == nil {
		return readConfig()
	} else {
		return models.DefaultConfig()
	}
}

func main() {
	cfg := loadConfig()

	if len(os.Args) > 1 {
		c := client.New(api.DefaultConfig(), cfg.Namespace)
		g := git.New(cfg)

		cli := NewCLI(c, g)
		cli.Run()
	} else {
		fmt.Println(Usage)
	}
}
