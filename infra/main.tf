data "oci_identity_tenancy" "tenancy" {
  tenancy_id = var.tenancy_id
}

resource "oci_identity_compartment" "freecloud" {
  compartment_id = data.oci_identity_tenancy.tenancy.id
  name           = var.compartment_name
  description    = var.compartment_description
  freeform_tags  = var.common_tags
}

module "base" {
  source         = "./modules/base"
  compartment_id = oci_identity_compartment.freecloud.id
  tags           = var.common_tags
}

module "vm" {
  count          = 1
  source         = "./modules/node"
  compartment_id = oci_identity_compartment.freecloud.id
  subnet_id      = module.base.subnet_id
  ssh_public_key = file("~/.ssh/id_ed25519.pub")
  shape = {
    name   = "VM.Standard.A1.Flex"
    config = {
      cpus   = 4
      memory = 24
    }
  }
  tags           = var.common_tags
}
