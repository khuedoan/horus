terraform {
  backend "remote" {
    organization = "khuedoan"

    workspaces {
      name = "freecloud"
    }
  }
}
