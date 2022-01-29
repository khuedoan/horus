.POSIX:

default: infra

.PHONY: infra
infra:
	make -C infra
