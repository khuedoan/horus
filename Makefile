.POSIX:
.PHONY: default infra cluster edit-vault test update

default: infra cluster

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

edit-vault:
	ansible-vault edit ./cluster/roles/global-secrets/vars/main.yml

test:
	cd test/e2e && go test

update:
	nix flake update
