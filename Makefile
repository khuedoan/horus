.POSIX:
.PHONY: *

default: infra cluster

infra:
	make -C infra

cluster:
	make -C cluster

update:
	nix flake update

test:
	cd test/e2e && go test
