OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)

## Tool Versions
# Set the OPM version to use.
OPM_VERSION ?= v1.56.0

all: render-all-catalogs

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## Location to install dependencies to
LOCALBIN := $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: opm
OPM_DL_URL = https://github.com/operator-framework/operator-registry/releases/download/$(OPM_VERSION)
OPM = $(LOCALBIN)/opm-$(OPM_VERSION)
opm: $(LOCALBIN) opm-policy ## Download opm locally if necessary.
ifeq (,$(wildcard $(OPM)))
	@{ \
	set -e ;\
	echo "Downloading OPM version $(OPM_VERSION)" ;\
	curl -sSLo $(OPM) $(OPM_DL_URL)/$(OS)-$(ARCH)-opm ;\
	chmod +x $(OPM) ;\
	}
endif

.PHONY: opm-policy
POLICY_FILE = $(HOME)/.config/containers/policy.json
opm-policy:
ifeq (,$(wildcard $(POLICY_FILE)))
	@{ \
	set -e ;\
	if ! test -f $(POLICY_FILE); then \
		echo "Creating policy.json file at $(POLICY_FILE)"; \
		mkdir -p $(dir $(POLICY_FILE)); \
		echo '{ "default": [{"type": "insecureAcceptAnything"}] }' > $(POLICY_FILE); \
	fi \
	}
endif

.PHONY: render-everest-operator-catalog
render-everest-operator-catalog: opm ## Render the everest-operator catalog.
	@echo "Rendering everest-operator catalog"
	$(OPM) alpha render-template basic -o yaml ./veneer/everest-operator.yaml > ./catalog/everest-operator/catalog.yaml

.PHONY: render-postgresql-operator-catalog
render-postgresql-operator-catalog: opm ## Render the postgresql-operator catalog.
	@echo "Rendering postgresql-operator catalog"
	$(OPM) alpha render-template semver -o yaml ./veneer/percona-postgresql-operator.yaml > ./catalog/percona-postgresql-operator/catalog.yaml

.PHONY: render-psmdb-operator-catalog
render-psmdb-operator-catalog: opm ## Render the psmdb-operator catalog.
	@echo "Rendering psmdb-operator catalog"
	$(OPM) alpha render-template semver -o yaml ./veneer/percona-server-mongodb-operator.yaml > ./catalog/percona-server-mongodb-operator/catalog.yaml

.PHONY: render-pxc-operator-catalog
render-pxc-operator-catalog: opm ## Render the pxc-operator catalog.
	@echo "Rendering pxc-operator catalog"
	$(OPM) alpha render-template semver -o yaml ./veneer/percona-xtradb-cluster-operator.yaml > ./catalog/percona-xtradb-cluster-operator/catalog.yaml

.PHONY: render-ps-operator-catalog
render-ps-operator-catalog: opm ## Render the ps-operator catalog.
	@echo "Rendering ps-operator catalog"
	$(OPM) alpha render-template semver -o yaml ./veneer/percona-server-mysql-operator.yaml > ./catalog/percona-server-mysql-operator/catalog.yaml

.PHONY: render-vm-operator-catalog
render-vm-operator-catalog: opm ## Render the vm-operator catalog.
	@echo "Rendering vm-operator catalog"
	$(OPM) alpha render-template semver -o yaml ./veneer/victoriametrics-operator.yaml > ./catalog/victoriametrics-operator/catalog.yaml

PHONY: render-all-catalogs
TARGET_DEPS := render-everest-operator-catalog render-postgresql-operator-catalog
TARGET_DEPS += render-psmdb-operator-catalog render-pxc-operator-catalog
TARGET_DEPS += render-ps-operator-catalog render-vm-operator-catalog
render-all-catalogs: $(TARGET_DEPS) ## Render the vm-operator catalog
