# dcdr
Decider Feature flags CLI

# Overview

Decider is a feature flag system built with the Consul KV Store. It supports both `percentile` and `boolean` flags for controlled infrastructure rollouts and releases. Changes within the configurable KV namespace are observed and distributed to other nodes in a cluster via Consul Template which generates a `JSON` file containing all current flags and their values.

# Features
```go
type FeatureType string

const (
	Percentile FeatureType = "percentile"
	Boolean    FeatureType = "boolean"
)

type Feature struct {
	FeatureType FeatureType `json:"feature_type"`
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Comment     string      `json:"comment"`
	UpdatedBy   string      `json:"updated_by"`
}
```

# Installation 

See releases.

### Install from source

```bash
git clone https://github.com/vsco/dcdr.git
cd dcdr
./script/install
```

# Usage 

```bash
$ dcdr

Usage:

        dcdr command [arguments]

The commands are:

        list        list all feature flags
        set         create or update a feature flag
        delete      delete a feature flag
        init        init the audit repo

Use "dcdr help [command]" for more information about a command.
```
