.POSIX:

default: infra bootstrap

.PHONY: infra
infra:
	make -C infra

.PHONY: bootstrap
bootstrap:
	make -C bootstrap
