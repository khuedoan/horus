terraform {
  backend "remote" {
    organization = "khuedoan"

    workspaces {
      name = "freeinfra"
    }
  }

  required_providers {
    oci = {
      source  = "hashicorp/oci"
      version = "4.40.0"
    }

    azurerm = {
      source = "hashicorp/azurerm"
      version = "2.73.0"
    }

    google = {
      source = "hashicorp/google"
      version = "3.81.0"
    }

    aws = {
      source = "hashicorp/aws"
      version = "3.55.0"
    }
  }
}
