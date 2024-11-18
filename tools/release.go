package main

import (
	"encoding/json"
	"fmt"
	goversion "github.com/hashicorp/go-version"
	"github.com/operator-framework/operator-registry/alpha/declcfg"
	"strings"
)

const (
	rcBundle      = "%s/perconalab/everest-operator-bundle:%s"
	releaseBundle = "%s/percona/everest-operator-bundle:%s"
)

type release struct {
	isRC            bool
	newVersion      goversion.Version
	previousVersion goversion.Version
	registry        string
}

func (r *release) image() string {
	version := r.newVersion.String()
	if r.isRC {
		return r.imageTemplate(rcBundle, version)
	}
	return r.imageTemplate(releaseBundle, version)
}

func (r *release) imageTemplate(bundleTemplate, tag string) string {
	return fmt.Sprintf(bundleTemplate, r.registry, tag)
}

func (r *release) create(currentVersionStr, newVersionStr, registry string) error {
	currentVersion, err := goversion.NewVersion(currentVersionStr)
	if err != nil {
		return fmt.Errorf("%s: can not parse current version: %w", currentVersionStr, err)
	}

	newVersion, err := goversion.NewVersion(newVersionStr)
	if err != nil {
		return fmt.Errorf("%s: can not parse new version: %w", newVersionStr, err)
	}

	if strings.Contains(newVersionStr, "rc") {
		r.isRC = true
	}
	r.newVersion = *newVersion
	r.previousVersion = *currentVersion
	r.registry = registry
	return nil
}

func (r *release) bundle() (*declcfg.Meta, error) {
	bundle := declcfg.Bundle{Schema: declcfg.SchemaBundle, Image: r.image()}
	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return nil, err
	}

	newEntry := declcfg.Meta{
		Schema: declcfg.SchemaBundle,
	}
	newEntry.Blob.UnmarshalJSON(bundleBytes)
	return &newEntry, nil
}
