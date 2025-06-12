.POSIX:
.PHONY: default compose infra apps test update

env ?= local
# TODO multiple clusters
export KUBECONFIG = $(shell pwd)/cluster/kubeconfig.yaml

default: infra

compose:
	docker compose up --build --detach

infra: compose
	# TODO multiple env
	temporal workflow start \
		--task-queue cloudlab \
		--type Infra \
		--input '{ "url": "https://github.com/khuedoan/cloudlab", "revision": "infra-rewrite", "oldRevision": "790763a8166e306f34559870c60e818505117e6b", "stack": "local" }'

apps:
	# TODO auto bootstrap
	kubectl apply --server-side=true --namespace argocd --filename apps/

platform:
	# TODO auto bootstrap
	kubectl apply --server-side=true --namespace argocd --filename platform/

test:
	cd controller && go test ./...
	cd test/e2e && go test

fmt:
	yamlfmt \
		--exclude infra/modules/cluster/roles/secrets/vars/main.yml \
		--exclude infra/*/secrets.yaml \
		.
	terragrunt hcl format
	tofu fmt -recursive
	cd controller && go fmt ./...
	cd test/e2e && go fmt ./...

update:
	nix flake update
