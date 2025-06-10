remote_state {
  backend = "s3"
  config = {
    endpoints = {
      s3 = "http://localhost:9000"
    }
    bucket                             = "tfstate"
    key                                = "${path_relative_to_include()}/terraform.tfstate"
    region                             = "eu-west-1"
    encrypt                            = false
    disable_aws_client_checksums       = true
    skip_bucket_ssencryption           = true
    skip_bucket_public_access_blocking = true
    skip_bucket_enforced_tls           = true
    skip_bucket_root_access            = true
    skip_credentials_validation        = true
    force_path_style                   = true
  }
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
}
