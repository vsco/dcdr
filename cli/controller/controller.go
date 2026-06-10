package controller

import (
	"encoding/json"
	"os"

	"errors"

	"io/ioutil"

	"path"

	"github.com/tucnak/climax"
	"github.com/vsco/dcdr/cli/api"
	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/cli/ui"
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
	"github.com/vsco/dcdr/server"
)

const filePerms = 0775

var (
	errInvalidType           = errors.New("invalid -type. use boolean, percentile, or string")
	errTypeRequiredForString = errors.New("-type=string is required for non-numeric, non-boolean values")
	errInvalidBool           = errors.New("invalid -value for boolean. use -value=[true|false]")
	errInvalidPercentile     = errors.New("invalid -value for percentile. must be a number")
	errEmptyString           = errors.New("invalid -value for string. must not be empty or whitespace-only")
	errInvalidRange          = errors.New("invalid -value for percentile. use -value=[0.0-1.0]")
	errNameRequired          = errors.New("-name is required")
)

// Controller handler for CLI commands
type Controller struct {
	Config *config.Config
	Client api.ClientIFace
}

// New creates a `Controller`
func New(cfg *config.Config, kv api.ClientIFace) (cc *Controller) {
	cc = &Controller{
		Config: cfg,
		Client: kv,
	}

	return
}

func (cc *Controller) List(ctx climax.Context) int {
	pf, _ := ctx.Get("prefix")
	scope, _ := ctx.Get("scope")

	if pf != "" && scope == "" {
		scope = models.DefaultScope
	}

	features, err := cc.Client.List(pf, scope)

	if err != nil {
		printer.SayErr("%v", err)
		return 1
	}

	if len(features) == 0 {
		printer.Say("no feature flags found in namespace: %s",
			cc.Client.Namespace())
		return 1
	}

	ui.New().DrawFeatures(features)

	return 0
}

func (cc *Controller) Set(ctx climax.Context) int {
	ft, err := cc.ParseContext(ctx)

	if err != nil {
		printer.SayErr("parse error: %v", err)
		return 1
	}

	err = cc.Client.Set(ft)

	if err != nil {
		printer.SayErr("set error: %v", err)
		return 1
	}

	printer.Say("set flag '%s'", ft.ScopedKey())

	return cc.CommitFeatures(ft, false)
}

func (cc *Controller) Delete(ctx climax.Context) int {
	name, _ := ctx.Get("name")
	scope, _ := ctx.Get("scope")

	if name == "" {
		printer.Say("name is required")
		return 1
	}

	if scope == "" {
		scope = models.DefaultScope
	}

	err := cc.Client.Delete(name, scope)

	if err != nil {
		printer.SayErr("%v", err)
		return 1
	}

	printer.Say("deleted flag %s/%s/%s",
		cc.Config.Namespace, scope, name)

	ft := &models.Feature{
		Key:       name,
		Scope:     scope,
		UpdatedBy: cc.Config.Username,
	}

	return cc.CommitFeatures(ft, true)
}

func (cc *Controller) CommitFeatures(ft *models.Feature, deleted bool) int {
	if cc.Config.GitEnabled() {
		printer.Say("committing changes")
		err := cc.Client.Commit(ft, deleted)

		if err != nil {
			printer.SayErr("%v", err)
			return 1
		}

		sha, err := cc.Client.UpdateCurrentSHA()
		printer.Say("set info/current_sha: %s", sha)

		if err != nil {
			printer.SayErr("%v", err)
			return 1
		}

		if cc.Config.PushEnabled() {
			printer.Say("pushing commit to origin")
			err = cc.Client.Push()

			if err != nil {
				printer.SayErr("%v", err)
				return 1
			}
		}

	}

	return 0
}

func (cc *Controller) Init(ctx climax.Context) int {
	if _, err := os.Stat(config.Path()); os.IsNotExist(err) {
		err = os.MkdirAll(path.Dir(config.Path()), filePerms)

		printer.Say("creating %s", path.Dir(config.Path()))

		if err != nil {
			printer.SayErr("could not create config directory: %v", err)
			return 1
		}

		err = ioutil.WriteFile(config.Path(), config.ExampleConfig, filePerms)
		printer.Say("%s not found. creating example config", config.Path())

		if err != nil {
			printer.SayErr("could not write config.hcl %v", err)
			return 1
		}
	}

	if !cc.Config.GitEnabled() {
		printer.Say("no repository configured. skipping")
		return 0
	}

	create := ctx.Is("create")

	err := cc.Client.InitRepo(create)

	if err != nil {
		printer.SayErr("%v", err)
		return 1
	}

	if create {
		printer.Say("initialized new repo in %s and pushed to %s",
			cc.Config.Git.RepoPath, cc.Config.Git.RepoURL)
	} else {
		printer.Say("cloned %s into %s",
			cc.Config.Git.RepoURL, cc.Config.Git.RepoPath)
	}

	return 0
}

func (cc *Controller) Import(ctx climax.Context) int {
	bts, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		printer.SayErr("%v", err)
		return 1
	}

	var kvs map[string]interface{}
	err = json.Unmarshal(bts, &kvs)

	if err != nil {
		printer.SayErr("%v", err)
		return 1
	}

	scope, _ := ctx.Get("scope")

	if scope == "" {
		scope = models.DefaultScope
	}

	for k, v := range kvs {
		f := models.NewFeature(k, v, "", "", scope, cc.Config.Namespace)
		err = cc.Client.Set(f)

		if err != nil {
			printer.SayErr("%v", err)
			return 1
		}

		printer.Say("set %s to %+v", k, v)
	}

	return 1
}

func (cc *Controller) Info(ctx climax.Context) int {

	ui.New().DrawConfig(cc.Config)

	return 0
}

func (cc *Controller) Serve(ctx climax.Context) int {
	c, err := client.New(cc.Config)

	if err != nil {
		printer.LogErrf("%v", err)
	}

	s := server.New(cc.Config, c)

	printer.Logf("pid: %d serving %s on %s", os.Getpid(),
		cc.Config.Server.Endpoint, cc.Config.Server.Host)

	err = s.Serve()

	if err != nil {
		printer.LogErrf("%v", err)
		return 1
	}

	return 0
}

func (cc *Controller) Watch(ctx climax.Context) int {
	printer.Logf("watching namespace: %s", cc.Config.Namespace)

	cc.Client.Watch()

	return 0
}

func (cc *Controller) ParseContext(ctx climax.Context) (*models.Feature, error) {
	name, _ := ctx.Get("name")
	val, _ := ctx.Get("value")
	typ, _ := ctx.Get("type")
	cmt, _ := ctx.Get("comment")
	scp, _ := ctx.Get("scope")

	if name == "" {
		return nil, errNameRequired
	}

	var v interface{}
	var ft models.FeatureType

	if val != "" {
		var err error
		v, ft, err = parseValue(val, typ)

		if err != nil {
			return nil, err
		}
	}

	f := models.NewFeature(name, v, cmt, cc.Config.Username, scp, cc.Config.Namespace)
	f.FeatureType = ft

	return f, nil
}

// parseValue resolves the value and feature type for a `set` command.
//   - When -type is provided, the value is parsed strictly for that type.
//   - When -type is omitted, the type is inferred: boolean and percentile are
//     accepted, but an inferred string is rejected (see below).
func parseValue(val string, typ string) (interface{}, models.FeatureType, error) {
	var v interface{}
	var ft models.FeatureType

	if typ != "" {
		var ok bool
		ft, ok = models.ParseFeatureType(typ)

		if !ok {
			return nil, models.Invalid, errInvalidType
		}

		var err error
		if v, err = models.ParseValueForType(val, ft); err != nil {
			switch ft {
			case models.Boolean:
				return nil, models.Invalid, errInvalidBool
			case models.Percentile:
				return nil, models.Invalid, errInvalidPercentile
			case models.String:
				return nil, models.Invalid, errEmptyString
			default:
				return nil, models.Invalid, err
			}
		}
	} else {
		v, ft = models.ParseValueAndFeatureType(val)

		// An inferred string is ambiguous (e.g. a typo'd bool/number), so
		// require the caller to opt in explicitly with -type=string.
		if ft == models.String {
			return nil, models.Invalid, errTypeRequiredForString
		}
	}

	// shared validation for both the explicit and inferred paths
	if ft == models.Percentile {
		if err := validatePercentile(v); err != nil {
			return nil, models.Invalid, err
		}
	}

	return v, ft, nil
}

// validatePercentile ensures a percentile value falls within the 0.0-1.0 range.
func validatePercentile(v interface{}) error {
	if f := v.(float64); f > 1.0 || f < 0 {
		return errInvalidRange
	}

	return nil
}
