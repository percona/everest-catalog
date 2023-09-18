
## Catalog Directory

This directory serves to contain individual operator catalogs.  Each operator will be contained in 
a sub-directory named after it.  
For example, for the operator catalog for 'OperatorA', we would anticipate the following directory structure:

```tree
catalog
├── .indexignore          (to make `opm validate` ignore README.md)
├── README.md
├── OperatorA
│   ├── .indexignore      (to make `opm validate` ignore OWNERS)
│   ├── OWNERS            (owners of the operator catalog to review PRs)
│   ├── veneer.yaml       (Veneer to build catalog.yaml)
│   └── catalog.yaml      (catalog compiled from veneer)
```

## NB:  
Because all levels of the catalog hierarchy need to be able to pass an opm validation attempt, it may be necessary to include `.indexignore` files to exclude non-catalog files from the attempt. 

[See the OpenShift documentation for more information on these files.](https://docs.openshift.com/container-platform/4.10/operators/understanding/olm-packaging-format.html)


## Add or Update operator catalog

Install latest `opm` tool: https://github.com/operator-framework/operator-registry/releases/latest

To add new operator - create `veneer.yaml`:

```sh
mkdir catalog/percona-xtradb-cluster-operator ; cd catalog/percona-xtradb-cluster-operator
cat << EOF >> veneer.yaml
Schema: olm.semver
GenerateMajorChannels: true
GenerateMinorChannels: false
Stable:
  Bundles:
  - Image: quay.io/operatorhubio/percona-xtradb-cluster-operator:v1.10.0
EOF
```

To update operator - add new bundle to the existing `veneer`:

```yaml
Schema: olm.semver
GenerateMajorChannels: true
GenerateMinorChannels: false
Stable:
  Bundles:
  - Image: quay.io/operatorhubio/percona-xtradb-cluster-operator:v1.10.0
  - Image: docker.io/percona/percona-xtradb-cluster-operator:v1.11.0-bundle
```
example-operator
Use channel names Candidate, Fast, and Stable as stated in [Veneer specification](https://olm.operatorframework.io/docs/reference/veneers/#specification).

Generate operator catalog:

```
opm alpha render-veneer semver -o yaml < veneer.yaml > catalog.yaml
opm validate .
```

Create PR and get it merged (`veener.yaml` and `catalog.yaml`).
