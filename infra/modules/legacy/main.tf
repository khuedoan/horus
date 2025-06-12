module "base" {
  source                  = "../base"
  tenancy_id              = var.tenancy_ocid
  compartment_name        = var.compartment_name
  compartment_description = var.compartment_description
}

module "network" {
  source         = "../network"
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
  source         = "../instance"
  compartment_id = module.base.compartment_id
  display_name   = "cloudlab"
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

resource "local_file" "inventory" {
  filename        = "${path.root}/inventory.yml"
  file_permission = "0644"
  content = yamlencode({
    k3s = {
      hosts = {
        "${module.instance.public_ip}" = {
          ansible_user                 = "ubuntu"
          ansible_ssh_private_key_file = abspath(local_file.ssh_private_key.filename)
        }
      }
    }
  })
}
