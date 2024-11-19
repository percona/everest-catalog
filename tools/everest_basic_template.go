package main

import (
	"encoding/json"
	"fmt"
	"github.com/operator-framework/operator-registry/alpha/declcfg"
	"github.com/operator-framework/operator-registry/alpha/template/basic"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const versionPrefix = "everest-operator.v"

type EverestBasicTemplate struct {
	basic.BasicTemplate
}

func (t *EverestBasicTemplate) readFromFile(veneerFile string) error {
	yamlFile, err := os.ReadFile(veneerFile)
	if err != nil {
		return err
	}

	var yamlData interface{}
	err = yaml.Unmarshal(yamlFile, &yamlData)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(yamlData)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, t)
}

func (t *EverestBasicTemplate) toByteArray() ([]byte, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	var jsonInterface interface{}
	err = json.Unmarshal(bytes, &jsonInterface)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(jsonInterface)
}

func (t *EverestBasicTemplate) update(build build, channel string) error {
	for i, meta := range t.Entries {
		if meta.Schema == declcfg.SchemaChannel && meta.Name == channel {
			bytes, err := updateBlob(meta.Blob, build, channel)
			if err != nil {
				return fmt.Errorf("failed to update channel data %s", channel)
			}
			t.Entries[i].Blob.UnmarshalJSON(bytes)
			break
		}
	}

	return t.addBundle(build)
}

func (t *EverestBasicTemplate) currentVersion(channel string) string {
	for _, meta := range t.Entries {
		if meta.Schema == declcfg.SchemaChannel && meta.Name == channel {
			var ch declcfg.Channel
			err := json.Unmarshal(meta.Blob, &ch)
			if err != nil {
				return ""
			}
			return strings.TrimPrefix(ch.Entries[len(ch.Entries)-1].Name, versionPrefix)
		}
	}
	return ""
}

func updateBlob(channelBlob json.RawMessage, build build, channel string) ([]byte, error) {
	var ch declcfg.Channel
	err := json.Unmarshal(channelBlob, &ch)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling channel %s: %w", channel, err)
	}

	if len(ch.Entries) == 0 {
		return nil, fmt.Errorf("versions list is empty for %s channel", channel)
	}
	lastVersion := ch.Entries[len(ch.Entries)-1]
	// new version inherits the replaces from the last version and adds a new skip to the skips
	newVersion := declcfg.ChannelEntry{
		Name:     versionPrefix + build.newVersion.String(),
		Replaces: lastVersion.Replaces,
		Skips:    append(lastVersion.Skips, lastVersion.Name),
	}
	// the previous version shouldn't have any replaces and skips anymore
	ch.Entries[len(ch.Entries)-1].Replaces = ""
	ch.Entries[len(ch.Entries)-1].Skips = nil
	// add the new version to the list of versions
	ch.Entries = append(ch.Entries, newVersion)
	return json.Marshal(ch)
}

func (t *EverestBasicTemplate) addBundle(b build) error {
	for _, entry := range t.Entries {
		var bundle declcfg.Bundle
		err := json.Unmarshal(entry.Blob, &bundle)
		if err != nil {
			continue
		}
		if bundle.Image == b.image() {
			// if there is already such a bundle defined - do nothing
			// it might happen when two channels are being updated separately with the same version.
			return nil
		}
	}
	bundle, err := b.bundle()
	if err != nil {
		return fmt.Errorf("failed to create new bundle: %w", err)
	}

	t.Entries = append(t.Entries, bundle)
	return nil
}
