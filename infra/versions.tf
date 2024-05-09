terraform {
  required_version = "~> 1.5.0"

  backend "remote" {
    organization = "khuedoan"

    workspaces {
      name = "horus"
    }
  }

  required_providers {
    oci = {
      source  = "hashicorp/oci"
      version = "~> 5.41.0"
    }
  }
}

provider "oci" {
  tenancy_ocid = var.oracle_cloud.tenancy_ocid
  user_ocid    = var.oracle_cloud.user_ocid
  fingerprint  = var.oracle_cloud.fingerprint
  private_key  = var.oracle_cloud.private_key
  region       = var.oracle_cloud.region
}
