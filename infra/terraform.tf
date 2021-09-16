terraform {
  backend "remote" {
    organization = "khuedoan"

    workspaces {
      name = "freecloud"
    }
  }

  required_providers {
    oci = {
      source  = "hashicorp/oci"
      version = "4.43.0"
    }
  }
}
