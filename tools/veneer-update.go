package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"

	goversion "github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

func main() {
	var (
		newVersion     string
		currentVersion string
		channel        string
		veneerFile     string
	)

	rootCmd := &cobra.Command{
		Use:   "veneer-update",
		Short: "Prints an updated veneer file which adds a minor or a patch version to the entries list",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				b   []byte
				err error
			)

			b, err = updateVeneer(veneerFile, channel, currentVersion, newVersion)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(b))
		},
	}

	rootCmd.Flags().StringVar(&newVersion, "new-version", "", "New version (e.g. 0.10.0)")
	rootCmd.Flags().StringVar(&currentVersion, "current-version", "", "Current version (e.g. 0.9.5)")
	rootCmd.Flags().StringVar(&channel, "channel", "", "Channel to update (e.g. stable-v0)")
	rootCmd.Flags().StringVar(&veneerFile, "veneer-file", "", "Path to veneer file")

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

type Entry struct {
	Schema         string    `yaml:"schema"`
	DefaultChannel string    `yaml:"defaultChannel,omitempty"`
	Name           string    `yaml:"name,omitempty"`
	Versions       []Version `yaml:"entries,omitempty"`
	Package        string    `yaml:"package,omitempty"`
	Image          string    `yaml:"image,omitempty"`
}

type Template struct {
	Schema  string  `yaml:"schema"`
	Entries []Entry `yaml:"entries"`
}

type Version struct {
	Name     string   `yaml:"name"`
	Replaces string   `yaml:"replaces,omitempty"`
	Skips    []string `yaml:"skips,omitempty"`
}

const (
	olmChannelSchema = "olm.channel"
	olmBundleSchema  = "olm.bundle"

	versionPrefix = "everest-operator.v"

	rcBundle      = "docker.io/perconalab/everest-operator-bundle:"
	releaseBundle = "docker.io/percona/everest-operator-bundle:"
)

func updateVeneer(veneerFile, channel, currentVersionStr, newVersionStr string) ([]byte, error) {
	release, err := validateVersion(currentVersionStr, newVersionStr)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid version format: %w", newVersionStr, err)
	}

	source, err := os.ReadFile(veneerFile)
	if err != nil {
		return nil, fmt.Errorf("could not read veneer file: %w", err)
	}

	var template Template
	err = yaml.Unmarshal(source, &template)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return nil, err
	}

	for i, schema := range template.Entries {
		if schema.Schema == olmChannelSchema && schema.Name == channel {
			if len(schema.Versions) == 0 {
				return nil, fmt.Errorf("versions list is empty for %s channel", channel)
			}
			lastVersion := schema.Versions[len(schema.Versions)-1]
			// new version inherits the replaces from the last version and adds a new skip to the skips
			newVersion := Version{
				Name:     versionPrefix + release.newVersion.String(),
				Replaces: lastVersion.Replaces,
				Skips:    append(lastVersion.Skips, lastVersion.Name),
			}
			// the previous version shouldn't have any replaces and skips anymore
			schema.Versions[len(schema.Versions)-1].Replaces = ""
			schema.Versions[len(schema.Versions)-1].Skips = nil
			// add the new version to the list of versions
			schema.Versions = append(schema.Versions, newVersion)
			template.Entries[i] = schema
			break
		}
	}

	template.Entries = append(template.Entries, Entry{Schema: olmBundleSchema, Image: release.image()})
	return yaml.Marshal(template)
}

type release struct {
	isRC            bool
	newVersion      goversion.Version
	previousVersion goversion.Version
}

func (r release) image() string {
	version := r.newVersion.String()
	if r.isRC {
		return rcBundle + version
	}
	return releaseBundle + version
}

func validateVersion(currentVersionStr, newVersionStr string) (release, error) {
	result := release{}
	currentVersion, err := goversion.NewVersion(currentVersionStr)
	if err != nil {
		return result, fmt.Errorf("%s: can not parse current version: %w", currentVersionStr, err)
	}

	newVersion, err := goversion.NewVersion(newVersionStr)
	if err != nil {
		return result, fmt.Errorf("%s: can not parse new version: %w", newVersionStr, err)
	}

	if strings.Contains(newVersionStr, "rc") {
		result.isRC = true
	}
	result.newVersion = *newVersion
	result.previousVersion = *currentVersion
	return result, nil
}
