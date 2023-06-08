.POSIX:

default: apply

~/.terraform.d/credentials.tfrc.json:
	terraform login

~/.oci:
	oci setup config

.terraform.lock.hcl: ~/.terraform.d/credentials.tfrc.json ~/.oci versions.tf modules
	terraform init
	touch .terraform.lock.hcl

terraform.tfvars:
	cp terraform.tfvars.example terraform.tfvars
	nvim terraform.tfvars

plan: terraform.tfvars .terraform.lock.hcl
	terraform plan

apply: terraform.tfvars .terraform.lock.hcl
	terraform apply
