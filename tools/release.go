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

type build struct {
	isRelease       bool
	newVersion      goversion.Version
	previousVersion goversion.Version
	registry        string
}

func (b *build) image() string {
	version := b.newVersion.String()
	if b.isRelease {
		return b.imageTemplate(releaseBundle, version)
	}
	return b.imageTemplate(rcBundle, version)
}

func (b *build) imageTemplate(bundleTemplate, tag string) string {
	return fmt.Sprintf(bundleTemplate, b.registry, tag)
}

func (b *build) create(currentVersionStr, newVersionStr, registry string, testRepo bool) error {
	currentVersion, err := goversion.NewVersion(currentVersionStr)
	if err != nil {
		return fmt.Errorf("%s: can not parse current version: %w", currentVersionStr, err)
	}

	newVersion, err := goversion.NewVersion(newVersionStr)
	if err != nil {
		return fmt.Errorf("%s: can not parse new version: %w", newVersionStr, err)
	}

	b.newVersion = *newVersion
	b.previousVersion = *currentVersion
	b.registry = registry
	if !strings.Contains(newVersionStr, "rc") && !testRepo {
		b.isRelease = true
	}
	return nil
}

func (b *build) bundle() (*declcfg.Meta, error) {
	bundle := declcfg.Bundle{Schema: declcfg.SchemaBundle, Image: b.image()}
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
