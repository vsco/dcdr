package cli

import (
	"fmt"
	"strconv"

	"io/ioutil"
	"os"

	"encoding/json"

	"github.com/tucnak/climax"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/kv"
	"github.com/vsco/dcdr/models"
	"github.com/vsco/dcdr/repo"
	"github.com/vsco/dcdr/ui"
)

type Controller struct {
	Config *config.Config
	Store  *kv.Client
	Repo   *repo.Git
}

func NewController(cfg *config.Config, kv *kv.Client, repo *repo.Git) (cc *Controller) {
	cc = &Controller{
		Config: cfg,
		Store:  kv,
		Repo:   repo,
	}

	return
}

func (cc *Controller) List(ctx climax.Context) int {
	pf, _ := ctx.Get("prefix")
	features, err := cc.Store.List(pf)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	if len(features) == 0 {
		fmt.Printf("No feature flags found in namespace: %s.\n", cc.Store.Namespace)
		return 1
	}

	ui.New().DrawTable(features)

	return 0
}

func (cc *Controller) Set(ctx climax.Context) int {
	if err := cc.checkRepo(); err != nil {
		fmt.Println(err)
		return 1
	}

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

	existing, _ := cc.Store.Get(name)

	if existing != nil {
		ftc = existing.FeatureType

	} else {
		ftc = models.GetFeatureType(ft)
	}

	switch ftc {
	case models.Percentile:
		var v float64
		var err error

		if val == "" && existing != nil {
			v = existing.Value.(float64)
		} else {
			v, err = strconv.ParseFloat(val, 64)

			if err != nil {
				fmt.Println("invalid -value format. use -value=[0.0-1.0]")
				return 1
			}
		}

		cc.Store.SetPercentile(name, v, cmt, cc.Config.Username)
	case models.Boolean:
		var v bool
		var err error

		if val == "" && existing != nil {
			v = existing.Value.(bool)
		} else {
			v, err = strconv.ParseBool(val)

			if err != nil {
				fmt.Println("invalid -value format. use -value=[true,false]")
				return 1
			}
		}

		cc.Store.SetBoolean(name, v, cmt, cc.Config.Username)
	default:
		fmt.Printf("%q is not valid type.\n", ft)
		return 1
	}

	features, err := cc.Store.List("")

	if err != nil {
		fmt.Println(err)
		return 1
	}

	if cc.Config.UseGit() {
		if err := cc.Repo.Commit(features, msg); err != nil {
			fmt.Println(err)
			return 1
		}

		err := cc.updateCurrentSha()

		if err != nil {
			fmt.Println(err)
			return 1
		}
	}

	fmt.Printf("set flag '%s'\n", name)

	return 0
}

func (cc *Controller) Delete(ctx climax.Context) int {
	if err := cc.checkRepo(); err != nil {
		fmt.Println(err)
		return 1
	}

	name, _ := ctx.Get("name")
	if name == "" {
		fmt.Println("name is required")
		return 1
	}

	err := cc.Store.Delete(name)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	features, err := cc.Store.List("")

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

func (cc *Controller) Init(ctx climax.Context) int {
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

func (cc *Controller) Import(ctx climax.Context) int {
	bts, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	var kvs map[string]interface{}
	err = json.Unmarshal(bts, &kvs)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	for k, v := range kvs {
		switch v.(type) {
		case bool:
			cc.Store.SetBoolean(k, v.(bool), "", "")
		case float64:
			cc.Store.SetPercentile(k, v.(float64), "", "")
		default:
			fmt.Printf("skipped %s: unsupported type\n", k)
			continue
		}

		fmt.Printf("set %s to %+v\n", k, v)
	}

	return 1
}

func (cc *Controller) updateCurrentSha() error {
	sha, err := cc.Repo.CurrentSha()

	if err != nil {
		return err
	}

	err = cc.Store.SetCurrentSha(sha)

	return err
}

func (cc *Controller) checkRepo() error {
	if cc.Config.UseGit() && !cc.Repo.RepoExists() {
		return fmt.Errorf("%s does not exist. see `dcdr help init` for usage\n", cc.Config.Git.RepoPath)
	}

	return nil
}
