package main

import (
	"fmt"
	"os"

	"flag"

	"strconv"

	"github.com/hashicorp/consul/api"
	"github.com/vsco/decider-cli/client"
	"github.com/vsco/decider-cli/models"
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
		ft := set.String("type", "percentile", "the feature type [percentile,boolean,scalar]")
		val := set.String("value", "0.0", "the feature value")
		cmt := set.String("comment", "", "additional comment")

		set.Parse(os.Args[2:])

		ftc := models.GetFeatureType(*ft)

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
		case models.Scalar:
			f, err := strconv.ParseFloat(*val, 64)

			if err != nil {
				fmt.Println("invalid -value format. use -value=[0.0-1.0]")
				os.Exit(2)
			}

			c.client.SetScalar(*name, f, *cmt)
		default:
			fmt.Printf("%q is not valid type.\n", *ft)
			os.Exit(2)
		}

		features, err := c.client.List(*name)

		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		ui.New().DrawTable(features)

	case "delete":
		set := flag.NewFlagSet("delete", flag.ExitOnError)
		n := set.String("name", "", "the feature name")

		set.Parse(os.Args[2:])

		err := c.client.Delete(*n)

		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		fmt.Printf("Deleted feature '%s'.\n", *n)
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}

func main() {
	c := client.New(api.DefaultConfig(), "decider/features")

	cli := NewCLI(c)
	cli.Run()
}
