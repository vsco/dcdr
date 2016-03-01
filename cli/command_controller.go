package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"errors"

	"os/exec"

	"io/ioutil"

	"github.com/tucnak/climax"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/kv"
	"github.com/vsco/dcdr/models"
	"github.com/vsco/dcdr/ui"
)

var (
	InvalidPercentileFormat = errors.New("invalid -value format. use -value=[0.0-1.0]")
	InvalidBoolFormat       = errors.New("invalid -value format. use -value=[true,false]")
	InvalidFeatureType      = errors.New("invalid -type. use -type=[boolean|percentile]")
)

type Controller struct {
	Config *config.Config
	Client kv.ClientIFace
}

func NewController(cfg *config.Config, kv kv.ClientIFace) (cc *Controller) {
	cc = &Controller{
		Config: cfg,
		Client: kv,
	}

	return
}

func (cc *Controller) Watch(ctx climax.Context) int {
	cmd := exec.Command(
		"consul",
		"watch",
		"-type",
		"keyprefix",
		"-prefix",
		cc.Config.Namespace,
		"cat")

	pr, pw := io.Pipe()
	cmd.Stdout = pw
	b := &bytes.Buffer{}
	cmd.Stderr = b

	scanner := bufio.NewScanner(pr)
	scanner.Split(bufio.ScanLines)

	go func() {
		for scanner.Scan() {
			fts, err := models.ParseFeatures(scanner.Bytes())

			if err != nil {
				fmt.Println(err)
				return
			}

			bts, err := fts.ToJSON()

			if err != nil {
				log.Println(err)
				return
			}

			err = ioutil.WriteFile(cc.Config.FilePath, bts, 0644)

			if err != nil {
				log.Println(err)
				return
			}

			log.Printf("wrote changes to %s\n", cc.Config.FilePath)
		}

		if scanner.Err() != nil {
			fmt.Println(scanner.Err())
		}
	}()

	log.Printf("watching namespace: %s\n", cc.Config.Namespace)

	if err := cmd.Run(); err != nil {
		log.Println(err, b)
	}

	return 0
}

func (cc *Controller) List(ctx climax.Context) int {
	pf, _ := ctx.Get("prefix")
	scope, _ := ctx.Get("scope")

	if pf != "" && scope == "" {
		scope = models.DefaultScope
	}

	features, err := cc.Client.List(pf, scope)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	if len(features) == 0 {
		fmt.Printf("No feature flags found in namespace: %s.\n", cc.Client.Namespace())
		return 1
	}

	ui.New().DrawTable(features)

	return 0
}

func (cc *Controller) ParseContext(ctx climax.Context) (*kv.SetRequest, error) {
	name, _ := ctx.Get("name")
	val, _ := ctx.Get("value")
	typ, _ := ctx.Get("type")
	cmt, _ := ctx.Get("comment")
	scp, _ := ctx.Get("scope")
	ft := models.GetFeatureType(typ)

	var v interface{}
	var err error

	switch ft {
	case models.Percentile:
		v, err = strconv.ParseFloat(val, 64)

		if err != nil {
			return nil, InvalidPercentileFormat
		}
	case models.Boolean:
		v, err = strconv.ParseBool(val)

		if err != nil {
			return nil, InvalidBoolFormat
		}
	case models.Invalid:
		return nil, InvalidFeatureType
	}

	return &kv.SetRequest{
		Key:       name,
		Value:     v,
		Scope:     scp,
		Namespace: cc.Config.Namespace,
		Comment:   cmt,
		User:      cc.Config.Username,
	}, nil
}

func (cc *Controller) Set(ctx climax.Context) int {
	sr, err := cc.ParseContext(ctx)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	err = cc.Client.Set(sr)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Printf("set flag '%s'\n", sr.Key)

	return 0
}

func (cc *Controller) Delete(ctx climax.Context) int {
	name, _ := ctx.Get("name")
	scope, _ := ctx.Get("scope")

	if name == "" {
		fmt.Println("name is required")
		return 1
	}

	if scope == "" {
		scope = models.DefaultScope
	}

	err := cc.Client.Delete(name, scope)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Printf("deleted flag %s/%s\n", scope, name)

	return 0
}

func (cc *Controller) Init(ctx climax.Context) int {
	_, create := ctx.Get("create")

	err := cc.Client.InitRepo(create)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	if create {
		fmt.Printf("initialized new repo in %s and pushed to %s", cc.Config.Git.RepoPath, cc.Config.Git.RepoURL)
	} else {
		fmt.Printf("cloned %s into %s", cc.Config.Git.RepoURL, cc.Config.Git.RepoPath)
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

	scope, _ := ctx.Get("scope")

	if scope == "" {
		scope = models.DefaultScope
	}

	for k, v := range kvs {
		sr := &kv.SetRequest{
			Key:       k,
			Value:     v,
			Namespace: cc.Config.Namespace,
			Scope:     scope,
		}

		err = cc.Client.Set(sr)

		if err != nil {
			fmt.Println(err)
			return 1
		}

		fmt.Printf("set %s to %+v\n", k, v)
	}

	return 1
}
