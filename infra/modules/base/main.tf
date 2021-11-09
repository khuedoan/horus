locals {
  protocols = {
    all    = "all",
    icmp   = "1",
    icmpv6 = "58",
    tcp    = "6",
    udp    = "17"
  }
}

resource "oci_core_vcn" "vcn" {
  compartment_id = var.compartment_id
  display_name   = "vcn"
  cidr_blocks    = var.vcn_cidr_blocks
  freeform_tags  = var.tags
}

resource "oci_core_security_list" "base" {
  compartment_id = var.compartment_id
  display_name   = "base"
  vcn_id         = oci_core_vcn.vcn.id
  freeform_tags  = var.tags

  ingress_security_rules {
    description = "Wireguard"
    protocol    = local.protocols["udp"]
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
    protocol    = local.protocols["tcp"]
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
    protocol    = local.protocols["tcp"]
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

  # TODO tighten this security rule
  ingress_security_rules {
    description = "Kube"
    protocol    = local.protocols["all"]
    source      = "0.0.0.0/0"
    stateless   = false
  }
}

resource "oci_core_subnet" "subnet" {
  compartment_id = var.compartment_id
  display_name   = "subnet"
  cidr_block     = var.subnet_cidr_block
  route_table_id = oci_core_vcn.vcn.default_route_table_id
  vcn_id         = oci_core_vcn.vcn.id
  freeform_tags  = var.tags

  security_list_ids = [
    oci_core_vcn.vcn.default_security_list_id,
    oci_core_security_list.base.id
  ]
}

resource "oci_core_internet_gateway" "internet_gateway" {
  compartment_id = var.compartment_id
  display_name   = "internet-gateway"
  vcn_id         = oci_core_vcn.vcn.id
  freeform_tags  = var.tags
}

resource "oci_core_default_route_table" "default_route_table" {
  manage_default_resource_id = oci_core_vcn.vcn.default_route_table_id
  route_rules {
    destination       = "0.0.0.0/0"
    destination_type  = "CIDR_BLOCK"
    network_entity_id = oci_core_internet_gateway.internet_gateway.id
  }
  freeform_tags = var.tags
}
