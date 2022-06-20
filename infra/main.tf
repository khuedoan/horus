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

resource "tls_private_key" "ssh" {
  algorithm = "ED25519"
}

resource "local_file" "ssh_private_key" {
  content         = tls_private_key.ssh.private_key_openssh
  filename        = "${path.root}/private.pem"
  file_permission = "0600"
}

module "instance" {
  source         = "./modules/instance"
  compartment_id = module.base.compartment_id
  display_name   = "horus"
  subnet_id      = module.network.subnet_id
  ssh_public_key = tls_private_key.ssh.public_key_openssh
  shape = {
    name = "VM.Standard.A1.Flex"
    config = {
      cpus   = 4
      memory = 24
    }
  }
}
