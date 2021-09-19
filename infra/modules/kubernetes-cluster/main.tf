data "oci_identity_availability_domains" "availability_domains" {
  compartment_id = var.compartment_id
}

data "oci_core_images" "ubuntu" {
  compartment_id = var.compartment_id
  operating_system = var.image_operating_system
  operating_system_version = var.image_operating_system_version
}

resource "oci_core_instance_configuration" "k3s_node" {
  compartment_id = var.compartment_id
}

resource "oci_core_instance_pool" "agent_pool" {
  compartment_id            = var.compartment_id
  instance_configuration_id = oci_core_instance_configuration.k3s_node.id
  size = var.agent_count

  placement_configurations {
    availability_domain = data.oci_identity_availability_domains.availability_domains.availability_domains[0].name
    primary_subnet_id   = var.subnet_id
  }
}
