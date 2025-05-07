.POSIX:
.PHONY: default infra cluster system platform apps secrets edit-secrets test update

# TODO multiple clusters
export KUBECONFIG = $(shell pwd)/cluster/kubeconfig.yaml

default: infra cluster system platform apps

~/.terraform.d/credentials.tfrc.json:
	# https://search.opentofu.org/provider/opentofu/tfe
	tofu login app.terraform.io

infra/terraform.tfvars:
	cp infra/terraform.tfvars.example infra/terraform.tfvars
	nvim infra/terraform.tfvars

infra: ~/.terraform.d/credentials.tfrc.json infra/terraform.tfvars
	cd infra \
		&& tofu init \
		&& tofu apply

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
	cd test/e2e && go fmt ./...

update:
	nix flake update
