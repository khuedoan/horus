.POSIX:
.PHONY: default compose infra platform apps test update

env ?= local

default: infra platform

compose:
	docker compose up --build --detach

infra: compose
	# TODO multiple env
	temporal workflow start \
		--task-queue cloudlab \
		--type Infra \
		--input '{ "url": "https://github.com/khuedoan/cloudlab", "revision": "master", "oldRevision": "790763a8166e306f34559870c60e818505117e6b", "stack": "local" }'

platform:
	# TODO multiple env
	temporal workflow start \
		--task-queue cloudlab \
		--type Platform \
		--input '{ "url": "https://github.com/khuedoan/cloudlab", "revision": "master", "registry": "registry.127.0.0.1.sslip.io", "cluster": "local" }'

apps:
	# TODO multiple env
	temporal workflow start \
		--task-queue cloudlab \
		--type Apps \
		--input '{ "url": "https://github.com/khuedoan/cloudlab", "revision": "master", "registry": "registry.127.0.0.1.sslip.io", "cluster": "local" }'

test:
	cd controller && go test ./...
	cd test && go test

fmt:
	nixfmt flake.nix
	yamlfmt \
		--exclude infra/modules/cluster/roles/secrets/vars/main.yml \
		--exclude infra/*/secrets.yaml \
		.
	terragrunt hcl format
	tofu fmt -recursive
	cd controller && go fmt ./...
	cd infra/modules/tfstate && go fmt ./...
	cd test && go fmt ./...

update:
	nix flake update

clean:
	docker compose down --remove-orphans --volumes
	k3d cluster delete local
