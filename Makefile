.POSIX:
.PHONY: *

default: infra cluster

infra:
	make -C infra

cluster:
	make -C cluster
