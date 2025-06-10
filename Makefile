.POSIX:
.PHONY: default docker-compose infra cluster system platform apps secrets edit-secrets test update

env ?= local
# TODO multiple clusters
export KUBECONFIG = $(shell pwd)/cluster/kubeconfig.yaml

default: infra cluster system platform apps

docker-compose:
	docker compose up --build --detach

~/.terraform.d/credentials.tfrc.json:
	# https://search.opentofu.org/provider/opentofu/tfe
	tofu login app.terraform.io

infra:
	cd infra/${env} \
		&& AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin terragrunt apply --all

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
	cd test/e2e && go test

fmt:
	yamlfmt --exclude cluster/roles/secrets/vars/main.yml .
	terragrunt hcl format
	tofu fmt -recursive
	cd controller && go fmt ./...
	cd test/e2e && go fmt ./...

update:
	nix flake update
