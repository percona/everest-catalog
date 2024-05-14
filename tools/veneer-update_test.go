package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatch(t *testing.T) {
	b, err := createPatchVersionVeneer("./everest-operator_test.yaml", "fast-v0", "0.9.1", "0.9.2")
	require.NoError(t, err)
	require.Equal(t, patchExpected, string(b))
}

func TestMinor(t *testing.T) {
	b, err := createMinorVersionVeneer("./everest-operator_test.yaml", "fast-v0", "0.9.1", "0.10.0")
	require.NoError(t, err)
	require.Equal(t, minorExpected, string(b))
}

var patchExpected = `defaultChannel: stable-v0
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
`

var minorExpected = `defaultChannel: stable-v0
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
`
