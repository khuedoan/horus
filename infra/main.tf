module "base" {
  source                  = "./modules/base"
  tenancy_id              = var.tenancy_id
  compartment_name        = var.compartment_name
  compartment_description = var.compartment_description
}

module "network" {
  source         = "./modules/network"
  compartment_id = module.base.compartment_id
}

module "instance" {
  source         = "./modules/instance"
  compartment_id = module.base.compartment_id
  display_name   = "freecloud"
  subnet_id      = module.network.subnet_id
  ssh_public_key = file("~/.ssh/id_ed25519.pub")
  shape = {
    name = "VM.Standard.A1.Flex"
    config = {
      cpus   = 4
      memory = 24
    }
  }
}
