## DBaaS Catalog

Welcome! This repository is a DBaaS Platform OLM Catalog which include sets of operators to enable DBaaS on Kubernetes.

See [Design](docs/design/dbaas-catalog-design.md) for more details.

### Use

```
minikube start

operator-sdk olm install

kubectl delete catalogsource operatorhubio-catalog -n olm

kubectl apply -f https://raw.githubusercontent.com/percona/dbaas-catalog/main/dbaas-catalog.yaml
kubectl get catalogsource -n olm
kubectl get packagemanifest -n olm
#wait for packagemanifests to appear

cat <<EOF | kubectl apply -f -
kind: OperatorGroup
apiVersion: operators.coreos.com/v1
metadata:
  name: og-single
  namespace: default
spec:
  targetNamespaces:
  - default
EOF

cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: victoriametrics-operator
  namespace: default
spec:
  channel: stable-v0
  installPlanApproval: Automatic
  name: victoriametrics-operator
  source: dbaas-catalog
  sourceNamespace: olm
  startingCSV: victoriametrics-operator.v0.27.2
EOF

cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: percona-xtradb-cluster-operator
  namespace: default
spec:
  channel: stable-v1
  installPlanApproval: Automatic
  name: percona-xtradb-cluster-operator
  source: dbaas-catalog
  sourceNamespace: olm
  startingCSV: percona-xtradb-cluster-operator.v1.10.0
EOF

cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: percona-server-mongodb-operator
  namespace: default
spec:
  channel: stable-v1
  installPlanApproval: Automatic
  name: percona-server-mongodb-operator
  source: dbaas-catalog
  sourceNamespace: olm
  startingCSV: percona-server-mongodb-operator.v1
EOF

kubectl get sub -n default
kubectl get csv -n default
kubectl get deployment -n default

```
