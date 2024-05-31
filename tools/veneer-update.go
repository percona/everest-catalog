package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	goversion "github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const imageNameTemplate = "docker.io/perconalab/everest-operator-bundle:%s"

func main() {
	var (
		versionType    string
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
			switch strings.ToLower(versionType) {
			case "minor":
				b, err = createMinorVersionVeneer(veneerFile, channel, currentVersion, newVersion)
			case "patch":
				b, err = createPatchVersionVeneer(veneerFile, channel, currentVersion, newVersion)
			default:
				panic("version-type is required")
			}

			if err != nil {
				panic(err)
			}

			out := strings.TrimSuffix(string(b), "---\n")
			fmt.Println(out)
		},
	}

	rootCmd.Flags().StringVar(&versionType, "version-type", "", "New version type, either minor or patch")
	rootCmd.Flags().StringVar(&newVersion, "new-version", "", "New version (e.g. 0.10.0)")
	rootCmd.Flags().StringVar(&currentVersion, "current-version", "", "Current version (e.g. 0.9.5)")
	rootCmd.Flags().StringVar(&channel, "channel", "", "Channel to update (e.g. stable-v0)")
	rootCmd.Flags().StringVar(&veneerFile, "veneer-file", "", "Path to veneer file")

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func createPatchVersionVeneer(veneerFile, channel, currentVersion, newVersion string) ([]byte, error) {
	source, err := os.ReadFile(veneerFile)
	if err != nil {
		return nil, fmt.Errorf("could not read veneer file: %w", err)
	}

	ret, err := updateVeneerFile(source, "olm.channel", "name", channel, func(m map[string]any) (map[string]any, error) {
		entries, ok := m["entries"]
		if !ok {
			return nil, errors.New("could not find entries")
		}

		eSlice, ok := entries.([]any)
		if !ok {
			return nil, errors.New("could not assert to slice")
		}

		i, e, err := findEntry(eSlice, currentVersion)
		if err != nil {
			return nil, err
		}

		// Keep just the name config in the current version
		eSlice[i] = map[string]any{"name": e["name"]}

		// Create a copy of the current version and rename to new version
		newEntry := e
		newEntry["name"] = fmt.Sprintf("everest-operator.v%s", newVersion)

		// Add current version to skips
		skips, ok := newEntry["skips"]
		if !ok {
			skips = []string{}
		}

		skipsSlice, ok := skips.([]any)
		if !ok {
			return nil, errors.New("could not assert skips to slice")
		}

		newEntry["skips"] = append(skipsSlice, fmt.Sprintf("everest-operator.v%s", currentVersion))

		eSlice = append(eSlice, newEntry)
		m["entries"] = eSlice

		return m, nil
	})

	if err != nil {
		return nil, err
	}

	return addOLMBundle(ret, newVersion)
}

func createMinorVersionVeneer(veneerFile, channel, currentVersion, newVersion string) ([]byte, error) {
	source, err := os.ReadFile(veneerFile)
	if err != nil {
		return nil, fmt.Errorf("could not read veneer file: %w", err)
	}

	ret, err := updateVeneerFile(source, "olm.channel", "name", channel, func(m map[string]any) (map[string]any, error) {
		entries, ok := m["entries"]
		if !ok {
			return nil, errors.New("could not find entries")
		}

		eSlice, ok := entries.([]any)
		if !ok {
			return nil, errors.New("could not assert to slice")
		}

		v, err := goversion.NewVersion(currentVersion)
		if err != nil {
			return nil, errors.Join(err, errors.New("could not parse current version"))
		}

		seg := v.Segments()
		skips := make([]string, 0, 10)
		for _, es := range eSlice {
			e, ok := es.(map[string]any)
			if !ok {
				return nil, errors.New("could not assert to map")
			}

			nameA, ok := e["name"]
			if !ok {
				continue
			}

			name, ok := nameA.(string)
			if !ok {
				continue
			}

			if strings.HasPrefix(name, fmt.Sprintf("everest-operator.v%d.%d.", seg[0], seg[1])) {
				skips = append(skips, name)
			}
		}

		newEntry := map[string]any{
			"name":     fmt.Sprintf("everest-operator.v%s", newVersion),
			"replaces": fmt.Sprintf("everest-operator.v%d.%d.0", seg[0], seg[1]),
			"skips":    skips,
		}

		m["entries"] = append(eSlice, newEntry)

		return m, nil
	})

	if err != nil {
		return nil, err
	}

	return addOLMBundle(ret, newVersion)
}

func findEntry(entries []any, version string) (int, map[string]any, error) {
	for i, es := range entries {
		e, ok := es.(map[string]any)
		if !ok {
			return 0, nil, errors.New("could not assert to map")
		}

		// Find document by name
		name, ok := e["name"]
		if !ok {
			continue
		}
		if name != fmt.Sprintf("everest-operator.v%s", version) {
			continue
		}

		return i, e, nil
	}

	return 0, nil, errors.New("not found")
}

func updateVeneerFile(veneerFile []byte, schema, fieldName, fieldValue string, replace func(map[string]any) (map[string]any, error)) ([]byte, error) {
	d := yaml.NewDecoder(bytes.NewReader(veneerFile))
	docs := make([]map[string]any, 0, 10)
	for {
		var doc map[string]any
		if err := d.Decode(&doc); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("could not decode yaml: %w", err)
		}

		if s, ok := doc["schema"]; !ok || s != schema {
			docs = append(docs, doc)
			continue
		}

		if f, ok := doc[fieldName]; !ok || f != fieldValue {
			docs = append(docs, doc)
			continue
		}

		replacement, err := replace(doc)
		if err != nil {
			return nil, err
		}
		docs = append(docs, replacement)
	}

	b := &bytes.Buffer{}
	for _, d := range docs {
		db, err := yaml.Marshal(d)
		if err != nil {
			return nil, fmt.Errorf("could not marshal yaml: %w", err)
		}
		if _, err := b.Write(db); err != nil {
			return nil, fmt.Errorf("could not write yaml: %w", err)
		}
		if _, err := b.WriteString("---\n"); err != nil {
			return nil, fmt.Errorf("could not write doc separator: %w", err)
		}
	}

	return b.Bytes(), nil
}

func addOLMBundle(veneeerContents []byte, newVersion string) ([]byte, error) {
	ret := veneeerContents

	var (
		foundImage bool
		imageName  = fmt.Sprintf(imageNameTemplate, newVersion)
	)

	_, err := updateVeneerFile(veneeerContents, "olm.bundle", "image", imageName, func(m map[string]any) (map[string]any, error) {
		foundImage = true
		return m, nil
	})
	if err != nil {
		return nil, errors.Join(err, errors.New("could not find olm.bundle"))
	}

	if !foundImage {
		ret = append(ret, []byte(fmt.Sprintf(
			"image: %s\n"+
				"schema: olm.bundle\n"+
				"---\n", imageName,
		))...)
	}

	return ret, nil
}
