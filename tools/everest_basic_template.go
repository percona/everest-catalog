package main

import (
	"encoding/json"
	"fmt"
	"github.com/operator-framework/operator-registry/alpha/declcfg"
	"github.com/operator-framework/operator-registry/alpha/template/basic"
	"gopkg.in/yaml.v3"
	"os"
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

func (t *EverestBasicTemplate) update(release release, channel string) error {
	for i, meta := range t.Entries {
		if meta.Schema == declcfg.SchemaChannel && meta.Name == channel {
			bytes, err := updateBlob(meta.Blob, release, channel)
			if err != nil {
				return fmt.Errorf("failed to update channel data %s", channel)
			}
			t.Entries[i].Blob.UnmarshalJSON(bytes)
			break
		}
	}

	bundle, err := release.bundle()
	if err != nil {
		return fmt.Errorf("failed to create new bundle: %w", err)
	}
	t.Entries = append(t.Entries, bundle)
	return nil
}

func updateBlob(channelBlob json.RawMessage, release release, channel string) ([]byte, error) {
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
		Name:     versionPrefix + release.newVersion.String(),
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
