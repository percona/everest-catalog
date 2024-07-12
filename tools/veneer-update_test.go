package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatch(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		channel        string
		currentVersion string
		newVersion     string
		expectedOutput string
	}{
		{
			name:           "Non-rc patch release",
			channel:        "fast-v0",
			currentVersion: "0.9.1",
			newVersion:     "0.9.2",
			expectedOutput: nonRCPatchExpected,
		},
		{
			name:           "RC patch release over non-rc release",
			channel:        "fast-v0",
			currentVersion: "0.9.1",
			newVersion:     "0.9.2-rc1",
			expectedOutput: rcPatchOverNonRCExpected,
		},
		{
			name:           "RC patch version over RC",
			channel:        "fast-v1",
			currentVersion: "0.9.1-rc1",
			newVersion:     "0.9.1-rc2",
			expectedOutput: rcOverRcPatchExpected,
		},
		{
			name:           "Non-RC minor release",
			channel:        "fast-v0",
			currentVersion: "0.9.1",
			newVersion:     "0.10.0",
			expectedOutput: nonRCMinorExpected,
		},
		{
			name:           "Minor release's first RC",
			channel:        "fast-v0",
			currentVersion: "0.9.1",
			newVersion:     "0.10.0-rc1",
			expectedOutput: minorReleasesFirstRCExpected,
		},
		{
			name:           "Minor release stable channel",
			channel:        "stable-v0",
			currentVersion: "0.9.1",
			newVersion:     "0.10.0",
			expectedOutput: minorReleaseStableChannelExpected,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b, err := updateVeneer("./everest-operator_test.yaml", tt.channel, tt.currentVersion, tt.newVersion)
			require.NoError(t, err)
			actual := string(b)
			require.Equal(t, tt.expectedOutput, actual)
		})
	}
}

var nonRCPatchExpected = `entries:
    - defaultChannel: stable-v0
      name: everest-operator
      schema: olm.package
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
        - name: everest-operator.v0.9.2
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
            - everest-operator.v0.9.1
      name: fast-v0
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1-rc1
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
      name: fast-v1
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
          skips:
            - everest-operator.v0.9.0
      name: stable-v0
      package: everest-operator
      schema: olm.channel
    - image: docker.io/perconalab/everest-operator-bundle:0.0.0
      schema: olm.bundle
    - image: docker.io/percona/everest-operator-bundle:0.9.2
      schema: olm.bundle
schema: olm.template.basic
`

var rcPatchOverNonRCExpected = `entries:
    - defaultChannel: stable-v0
      name: everest-operator
      schema: olm.package
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
        - name: everest-operator.v0.9.2-rc1
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
            - everest-operator.v0.9.1
      name: fast-v0
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1-rc1
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
      name: fast-v1
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
          skips:
            - everest-operator.v0.9.0
      name: stable-v0
      package: everest-operator
      schema: olm.channel
    - image: docker.io/perconalab/everest-operator-bundle:0.0.0
      schema: olm.bundle
    - image: docker.io/perconalab/everest-operator-bundle:0.9.2-rc1
      schema: olm.bundle
schema: olm.template.basic
`

var rcOverRcPatchExpected = `entries:
    - defaultChannel: stable-v0
      name: everest-operator
      schema: olm.package
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
      name: fast-v0
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1-rc1
        - name: everest-operator.v0.9.1-rc2
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
            - everest-operator.v0.9.1-rc1
      name: fast-v1
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
          skips:
            - everest-operator.v0.9.0
      name: stable-v0
      package: everest-operator
      schema: olm.channel
    - image: docker.io/perconalab/everest-operator-bundle:0.0.0
      schema: olm.bundle
    - image: docker.io/perconalab/everest-operator-bundle:0.9.1-rc2
      schema: olm.bundle
schema: olm.template.basic
`

var nonRCMinorExpected = `entries:
    - defaultChannel: stable-v0
      name: everest-operator
      schema: olm.package
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
        - name: everest-operator.v0.10.0
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
            - everest-operator.v0.9.1
      name: fast-v0
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1-rc1
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
      name: fast-v1
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
          skips:
            - everest-operator.v0.9.0
      name: stable-v0
      package: everest-operator
      schema: olm.channel
    - image: docker.io/perconalab/everest-operator-bundle:0.0.0
      schema: olm.bundle
    - image: docker.io/percona/everest-operator-bundle:0.10.0
      schema: olm.bundle
schema: olm.template.basic
`

var minorReleasesFirstRCExpected = `entries:
    - defaultChannel: stable-v0
      name: everest-operator
      schema: olm.package
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
        - name: everest-operator.v0.10.0-rc1
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
            - everest-operator.v0.9.1
      name: fast-v0
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1-rc1
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
      name: fast-v1
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
          skips:
            - everest-operator.v0.9.0
      name: stable-v0
      package: everest-operator
      schema: olm.channel
    - image: docker.io/perconalab/everest-operator-bundle:0.0.0
      schema: olm.bundle
    - image: docker.io/perconalab/everest-operator-bundle:0.10.0-rc1
      schema: olm.bundle
schema: olm.template.basic
`

var minorReleaseStableChannelExpected = `entries:
    - defaultChannel: stable-v0
      name: everest-operator
      schema: olm.package
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
      name: fast-v0
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.0.0
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1-rc1
          replaces: everest-operator.v0.0.0
          skips:
            - everest-operator.v0.0.0
            - everest-operator.v0.9.0
      name: fast-v1
      package: everest-operator
      schema: olm.channel
    - entries:
        - name: everest-operator.v0.9.0
        - name: everest-operator.v0.9.1
        - name: everest-operator.v0.10.0
          skips:
            - everest-operator.v0.9.0
            - everest-operator.v0.9.1
      name: stable-v0
      package: everest-operator
      schema: olm.channel
    - image: docker.io/perconalab/everest-operator-bundle:0.0.0
      schema: olm.bundle
    - image: docker.io/percona/everest-operator-bundle:0.10.0
      schema: olm.bundle
schema: olm.template.basic
`
