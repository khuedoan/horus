terraform {
  required_version = "~> 1.4.0"

  backend "remote" {
    organization = "khuedoan"

    workspaces {
      name = "horus"
    }
  }

  required_providers {
    oci = {
      source  = "hashicorp/oci"
      version = "~> 4.61.0"
    }
  }
}
