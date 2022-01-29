.POSIX:

default: infra cluster

.PHONY: infra
infra:
	make -C infra

.PHONY: cluster
cluster:
	make -C cluster
