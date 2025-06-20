locals {
  secrets = yamldecode(sops_decrypt_file(find_in_parent_folders("secrets.yaml")))
  env     = "production"
}

generate "backend" {
  path      = "backend.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
terraform {
  backend "s3" {
    bucket = "tfstate-${local.env}"
    key    = "${path_relative_to_include()}/tfstate.json"
    region                      = "auto"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true
    use_path_style              = true
    access_key = "${local.secrets.cloudflare_tfstate_access_key}"
    secret_key = "${local.secrets.cloudflare_tfstate_secret_key}"
    endpoints = { s3 = "https://${local.secrets.cloudflare_account_id}.r2.cloudflarestorage.com" }
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
