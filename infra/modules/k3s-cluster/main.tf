# module "server_pool" {
#   count  = var.server_count
#   source = "../virtual-machine"
#   subnet_id = var.subnet_id
#   ssh_public_key = var.ssh_public_key
# }

# module "agent_pool" {
#   count  = var.agent_count
#   source = "../virtual-machine"
#   subnet_id = var.subnet_id
#   ssh_public_key = var.ssh_public_key
# }

resource "oci_core_instance_configuration" "test_instance_configuration" {
  compartment_id = var.compartment_id
}

resource "oci_core_instance_pool" "agent_pool" {
  #Required
  compartment_id            = var.compartment_id
  instance_configuration_id = oci_core_instance_configuration.test_instance_configuration.id
  placement_configurations {
    #Required
    availability_domain = "gHLA:US-SANJOSE-1-AD-1"
    primary_subnet_id   = var.subnet_id

  }
  size = var.agent_count
}
