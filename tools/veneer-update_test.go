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
			name:           "Stable patch version",
			channel:        "fast-v0",
			currentVersion: "0.9.1",
			newVersion:     "0.9.2",
			expectedOutput: patchStableExpected,
		},
		{
			name:           "RC patch version over stable",
			channel:        "fast-v0",
			currentVersion: "0.9.1",
			newVersion:     "0.9.2-rc1",
			expectedOutput: patchStableRCExpected,
		},
		{
			name:           "RC patch version over RC",
			channel:        "fast-v1",
			currentVersion: "0.9.1-rc1",
			newVersion:     "0.9.1-rc2",
			expectedOutput: patchRCExpected,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b, err := createPatchVersionVeneer(
				"./everest-operator_test.yaml", tt.channel, tt.currentVersion, tt.newVersion,
			)
			require.NoError(t, err)
			require.Equal(t, tt.expectedOutput, string(b))
		})
	}
}

func TestMinor(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		channel        string
		currentVersion string
		newVersion     string
		expectedOutput string
	}{
		{
			name:           "Stable minor version",
			channel:        "fast-v0",
			currentVersion: "0.9.1",
			newVersion:     "0.10.0",
			expectedOutput: minorStableExpected,
		},
		{
			name:           "Stable minor RC version",
			channel:        "fast-v0",
			currentVersion: "0.9.1",
			newVersion:     "0.10.0-rc1",
			expectedOutput: minorStableRCExpected,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b, err := createMinorVersionVeneer(
				"./everest-operator_test.yaml", tt.channel, tt.currentVersion, tt.newVersion,
			)
			require.NoError(t, err)
			require.Equal(t, tt.expectedOutput, string(b))
		})
	}
}

var patchStableExpected = `defaultChannel: stable-v0
name: everest-operator
schema: olm.package
---
entries:
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
---
entries:
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
---
entries:
    - name: everest-operator.v0.9.0
    - name: everest-operator.v0.9.1
      skips:
        - everest-operator.v0.9.0
name: stable-v0
package: everest-operator
schema: olm.channel
---
image: docker.io/perconalab/everest-operator-bundle:0.0.0
schema: olm.bundle
---
image: docker.io/perconalab/everest-operator-bundle:0.9.2
schema: olm.bundle
---
`

var patchStableRCExpected = `defaultChannel: stable-v0
name: everest-operator
schema: olm.package
---
entries:
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
---
entries:
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
---
entries:
    - name: everest-operator.v0.9.0
    - name: everest-operator.v0.9.1
      skips:
        - everest-operator.v0.9.0
name: stable-v0
package: everest-operator
schema: olm.channel
---
image: docker.io/perconalab/everest-operator-bundle:0.0.0
schema: olm.bundle
---
image: docker.io/perconalab/everest-operator-bundle:0.9.2-rc1
schema: olm.bundle
---
`

var patchRCExpected = `defaultChannel: stable-v0
name: everest-operator
schema: olm.package
---
entries:
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
---
entries:
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
---
entries:
    - name: everest-operator.v0.9.0
    - name: everest-operator.v0.9.1
      skips:
        - everest-operator.v0.9.0
name: stable-v0
package: everest-operator
schema: olm.channel
---
image: docker.io/perconalab/everest-operator-bundle:0.0.0
schema: olm.bundle
---
image: docker.io/perconalab/everest-operator-bundle:0.9.1-rc2
schema: olm.bundle
---
`

var minorStableExpected = `defaultChannel: stable-v0
name: everest-operator
schema: olm.package
---
entries:
    - name: everest-operator.v0.0.0
    - name: everest-operator.v0.9.0
    - name: everest-operator.v0.9.1
      replaces: everest-operator.v0.0.0
      skips:
        - everest-operator.v0.0.0
        - everest-operator.v0.9.0
    - name: everest-operator.v0.10.0
      replaces: everest-operator.v0.9.0
      skips:
        - everest-operator.v0.9.0
        - everest-operator.v0.9.1
name: fast-v0
package: everest-operator
schema: olm.channel
---
entries:
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
---
entries:
    - name: everest-operator.v0.9.0
    - name: everest-operator.v0.9.1
      skips:
        - everest-operator.v0.9.0
name: stable-v0
package: everest-operator
schema: olm.channel
---
image: docker.io/perconalab/everest-operator-bundle:0.0.0
schema: olm.bundle
---
image: docker.io/perconalab/everest-operator-bundle:0.10.0
schema: olm.bundle
---
`

var minorStableRCExpected = `defaultChannel: stable-v0
name: everest-operator
schema: olm.package
---
entries:
    - name: everest-operator.v0.0.0
    - name: everest-operator.v0.9.0
    - name: everest-operator.v0.9.1
      replaces: everest-operator.v0.0.0
      skips:
        - everest-operator.v0.0.0
        - everest-operator.v0.9.0
    - name: everest-operator.v0.10.0-rc1
      replaces: everest-operator.v0.9.0
      skips:
        - everest-operator.v0.9.0
        - everest-operator.v0.9.1
name: fast-v0
package: everest-operator
schema: olm.channel
---
entries:
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
---
entries:
    - name: everest-operator.v0.9.0
    - name: everest-operator.v0.9.1
      skips:
        - everest-operator.v0.9.0
name: stable-v0
package: everest-operator
schema: olm.channel
---
image: docker.io/perconalab/everest-operator-bundle:0.0.0
schema: olm.bundle
---
image: docker.io/perconalab/everest-operator-bundle:0.10.0-rc1
schema: olm.bundle
---
`
