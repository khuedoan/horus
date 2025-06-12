locals {
  secrets = yamldecode(sops_decrypt_file(find_in_parent_folders("secrets.yaml")))
}

# TODO split into multiple modules, and use a more flexible state backend
generate "backend" {
  path      = "backend.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
terraform {
  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "khuedoan"

    workspaces {
      name = "cloudlab"
    }
  }
}
EOF
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
provider "oci" {
  tenancy_ocid = "${local.secrets.oracle_tenancy_ocid}"
  user_ocid    = "${local.secrets.oracle_user_ocid}"
  fingerprint  = "${local.secrets.oracle_fingerprint}"
  private_key  = <<EOT
${local.secrets.oracle_private_key}
EOT
  region       = "${local.secrets.oracle_region}"
}
EOF
}
