terraform {
  required_version = "~> 1.0.0"

  backend "remote" {
    organization = "khuedoan"

    workspaces {
      name = "freecloud"
    }
  }

  required_providers {
    oci = {
      source  = "hashicorp/oci"
      version = "~> 4.44.0"
    }
  }
}
