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

module "kubernetes_cluster" {
  source         = "./modules/kubernetes-cluster"
  compartment_id = oci_identity_compartment.freecloud.id
  server_count   = 1 # TODO multi master with embedded etcd in the same pool
  agent_count    = 1
  subnet_id      = module.base.subnet_id
  tags           = var.common_tags
}
