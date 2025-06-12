.POSIX:
.PHONY: default compose infra cluster system platform apps secrets edit-secrets test update

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

cluster:
	cd cluster && ansible-playbook \
		--inventory inventory.yml \
		--ask-vault-pass \
		main.yml

system:
	kubectl apply --server-side=true --namespace argocd --filename system/

platform:
	kubectl apply --server-side=true --namespace argocd --filename platform/

apps:
	kubectl apply --server-side=true --namespace argocd --filename apps/

secrets:
	cd cluster && ansible-playbook \
		--inventory inventory.yml \
		--ask-vault-pass \
		--tags secrets \
		main.yml

edit-secrets:
	ansible-vault edit ./cluster/roles/secrets/vars/main.yml

test:
	cd controller && go test ./...
	cd test/e2e && go test

fmt:
	yamlfmt \
		--exclude cluster/roles/secrets/vars/main.yml \
		--exclude infra/*/secrets.yaml \
		.
	terragrunt hcl format
	tofu fmt -recursive
	cd controller && go fmt ./...
	cd test/e2e && go fmt ./...

update:
	nix flake update
