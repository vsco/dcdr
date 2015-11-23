package main

import (
	"fmt"
	"log"
	"os"

	"flag"

	"strconv"

	"github.com/hashicorp/consul/api"
	"github.com/vsco/decider-cli/client"
	"github.com/vsco/decider-cli/ui"
)

type CLI struct {
	action string
	client *client.Client
}

func NewCLI(client *client.Client) (c *CLI) {
	c = &CLI{
		action: os.Args[1],
		client: client,
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

		if err != nil {
			log.Fatal(err)
		}

		ui.New().DrawTable(features)

	case "set":
		set := flag.NewFlagSet("set", flag.ExitOnError)
		n := set.String("name", "", "the feature name")
		ft := set.String("type", "percentile", "the feature type")
		val := set.String("value", "0.0", "the feature value")
		cmt := set.String("comment", "", "additional comment")

		set.Parse(os.Args[2:])

		switch *ft {
		case "percentile":
			f, err := strconv.ParseFloat(*val, 64)

			if err != nil {
				log.Fatal("invalid -value format. use -value=[0.0-1.0]")
			}

			c.client.SetPercentile(*n, f, *cmt)
		case "boolean":
			f, err := strconv.ParseBool(*val)

			if err != nil {
				log.Fatal("invalid -value format. use -value=[true,false]")
			}

			c.client.SetBoolean(*n, f, *cmt)
		case "scalar":
			f, err := strconv.ParseFloat(*val, 64)

			if err != nil {
				log.Fatal("invalid -value format. use -value=[0.0-1.0]")
			}

			c.client.SetScalar(*n, f, *cmt)
		default:
			fmt.Printf("%q is not valid type.\n", *ft)
			os.Exit(2)
		}

		features, err := c.client.List(*n)

		if err != nil {
			log.Fatal(err)
		}

		ui.New().DrawTable(features)

	case "delete":
		set := flag.NewFlagSet("delete", flag.ExitOnError)
		n := set.String("name", "", "the feature name")

		set.Parse(os.Args[2:])

		err := c.client.Delete(*n)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Deleted feature '%s'.\n", *n)
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}

func main() {
	c := client.New(api.DefaultConfig(), "decider")

	cli := NewCLI(c)
	cli.Run()
}
