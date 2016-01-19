package cli

import (
	"fmt"
	"strconv"

	"github.com/tucnak/climax"
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/git"
	"github.com/vsco/dcdr/models"
	"github.com/vsco/dcdr/ui"
)

type CommandController struct {
	Config *models.Config
	Client *client.Client
	Repo   *git.Git
}

func NewCommandController(cfg *models.Config, cl *client.Client, repo *git.Git) (cc *CommandController) {
	cc = &CommandController{
		Config: cfg,
		Client: cl,
		Repo:   repo,
	}

	return
}

func (cc *CommandController) List(ctx climax.Context) int {
	pf, _ := ctx.Get("prefix")
	features, err := cc.Client.List(pf)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	if len(features) == 0 {
		fmt.Printf("No feature flags found in namespace: %s.\n", cc.Client.Namespace)
		return 1
	}

	ui.New().DrawTable(features)

	return 0
}

func (cc *CommandController) Set(ctx climax.Context) int {
	name, _ := ctx.Get("name")
	val, _ := ctx.Get("value")
	ft, _ := ctx.Get("type")
	cmt, _ := ctx.Get("comment")

	msg := fmt.Sprintf("set %s to %s", name, val)

	if name == "" {
		fmt.Println("name is required")
		return 1
	}

	var ftc models.FeatureType

	existing, _ := cc.Client.Get(name)

	if existing != nil {
		ftc = existing.FeatureType
	} else {
		ftc = models.GetFeatureType(ft)
	}

	switch ftc {
	case models.Percentile:
		f, err := strconv.ParseFloat(val, 64)

		if err != nil {
			fmt.Println("invalid -value format. use -value=[0.0-1.0]")
			return 1
		}

		cc.Client.SetPercentile(name, f, cmt)
	case models.Boolean:
		f, err := strconv.ParseBool(val)

		if err != nil {
			fmt.Println("invalid -value format. use -value=[true,false]")
			return 1
		}

		cc.Client.SetBoolean(name, f, cmt)
	default:
		fmt.Printf("%q is not valid type.\n", ft)
		return 1
	}

	features, err := cc.Client.List("")

	if err != nil {
		fmt.Println(err)
		return 1
	}

	if err := cc.Repo.Commit(features, msg); err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Printf("set %s to %s\n", name, val)

	return 0
}

func (cc *CommandController) Delete(ctx climax.Context) int {
	name, _ := ctx.Get("name")
	if name == "" {
		fmt.Println("name is required")
		return 1
	}

	err := cc.Client.Delete(name)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	features, err := cc.Client.List("")

	if err != nil {
		fmt.Println(err)
		return 1
	}

	msg := fmt.Sprintf("deleted %s", name)
	if err := cc.Repo.Commit(features, msg); err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Printf("deleted flag '%s'\n", name)

	return 0
}

func (cc *CommandController) Init(ctx climax.Context) int {
	_, create := ctx.Get("create")

	if create {
		if err := cc.Repo.Create(); err != nil {
			fmt.Println(err)
			return 1
		} else {
			fmt.Printf("initialized new repo in %s and pushed to %s", cc.Config.Git.RepoPath, cc.Config.Git.RepoURL)
			return 0
		}
	}

	if err := cc.Repo.Clone(); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}
