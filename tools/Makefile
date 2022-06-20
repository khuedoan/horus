.POSIX:

TAG = horus-tools
WORKDIR = /usr/src/horus

default: build run

build:
	@docker build . --tag ${TAG}

run:
	@docker run \
		--rm \
		--interactive \
		--tty \
		--volume "$(shell pwd)/..:${WORKDIR}" \
		--volume "${HOME}/.oci:/root/.oci" \
		--volume "${HOME}/.terraform.d:/root/.terraform.d" \
		--workdir "${WORKDIR}" \
		${TAG}
