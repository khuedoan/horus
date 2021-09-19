resource "oci_core_vcn" "vcn" {
  cidr_blocks    = var.vcn_cidr_blocks
  compartment_id = var.compartment_id
  freeform_tags  = var.tags
}

resource "oci_core_security_list" "security_list" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.vcn.id
  freeform_tags  = var.tags

  ingress_security_rules {
    description = "Wireguard"
    protocol    = "17" # UDP
    source      = "0.0.0.0/0"
    stateless   = false

    udp_options {
      source_port_range {
        min = 1
        max = 65535
      }

      min = 51820
      max = 51820
    }
  }

  ingress_security_rules {
    description = "HTTP"
    protocol    = "6" # TCP
    source      = "0.0.0.0/0"
    stateless   = false

    tcp_options {
      source_port_range {
        min = 1
        max = 65535
      }

      min = 80
      max = 80
    }
  }

  ingress_security_rules {
    description = "HTTPS"
    protocol    = "6" # TCP
    source      = "0.0.0.0/0"
    stateless   = false

    tcp_options {
      source_port_range {
        min = 1
        max = 65535
      }

      min = 443
      max = 443
    }
  }
}

resource "oci_core_subnet" "subnet" {
  cidr_block     = var.subnet_cidr_block
  compartment_id = var.compartment_id
  route_table_id = oci_core_vcn.vcn.default_route_table_id
  vcn_id         = oci_core_vcn.vcn.id
  freeform_tags  = var.tags

  security_list_ids = [
    oci_core_vcn.vcn.default_security_list_id,
    oci_core_security_list.security_list.id
  ]
}

resource "oci_core_internet_gateway" "internet_gateway" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.vcn.id
  freeform_tags  = var.tags
}

resource "oci_core_default_route_table" "default_route_table" {
  route_rules {
    destination       = "0.0.0.0/0"
    destination_type  = "CIDR_BLOCK"
    network_entity_id = oci_core_internet_gateway.internet_gateway.id
  }
  manage_default_resource_id = oci_core_vcn.vcn.default_route_table_id
  freeform_tags  = var.tags
}
