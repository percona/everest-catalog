package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var (
		newVersion string
		channel    string
		veneerFile string
		registry   string
	)

	rootCmd := &cobra.Command{
		Use:   "veneer-update",
		Short: "Prints an updated veneer file which adds a minor or a patch version to the entries list",
		Run: func(cmd *cobra.Command, args []string) {
			err := updateVeneer(veneerFile, channel, newVersion, registry)
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd.Flags().StringVar(&newVersion, "new-version", "", "New version (e.g. 0.10.0)")
	rootCmd.Flags().StringVar(&channel, "channel", "", "Channel to update (e.g. stable-v0)")
	rootCmd.Flags().StringVar(&veneerFile, "veneer-file", "", "Path to veneer file")
	rootCmd.Flags().StringVar(&registry, "registry", "docker.io", "Registry to use")

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func updateVeneer(veneerFile, channel, newVersionStr, registry string) error {
	b, err := updateVeneerImpl(veneerFile, channel, newVersionStr, registry)
	if err != nil {
		return err
	}

	return os.WriteFile(veneerFile, b, 0644)
}

func updateVeneerImpl(veneerFile, channel, newVersionStr, registry string) ([]byte, error) {
	var t EverestBasicTemplate
	err := t.readFromFile(veneerFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	var r release
	err = r.create(t.currentVersion(channel), newVersionStr, registry)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid version format: %w", newVersionStr, err)
	}

	err = t.update(r, channel)
	if err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	b, err := t.toByteArray()
	if err != nil {
		return nil, fmt.Errorf("failed to convert data: %w", err)
	}

	return b, nil
}
