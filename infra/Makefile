.POSIX:

default: apply

~/.terraform.d/credentials.tfrc.json:
	terraform login

.terraform.lock.hcl: ~/.terraform.d/credentials.tfrc.json versions.tf modules
	tofu init
	touch .terraform.lock.hcl

terraform.tfvars:
	cp terraform.tfvars.example terraform.tfvars
	nvim terraform.tfvars

plan: terraform.tfvars .terraform.lock.hcl
	tofu plan

apply: terraform.tfvars .terraform.lock.hcl
	tofu apply
