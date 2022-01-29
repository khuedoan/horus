.POSIX:

default: infra cluster

.PHONY: infra
infra:
	make -C infra

.PHONY: cluster
cluster:
	make -C cluster

.PHONY: apps
apps:
	kubectl --kubeconfig=${PWD}/cluster/kubeconfig.yaml \
		apply \
		--filename apps
