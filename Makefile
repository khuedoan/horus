.POSIX:

default: infra config

.PHONY: infra
infra:
	make -C infra

.PHONY: config
config:
	make -C config
