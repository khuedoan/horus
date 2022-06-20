.POSIX:
.PHONY: *

default: infra cluster

tools:
	make -C tools

infra:
	make -C infra

cluster:
	make -C cluster
