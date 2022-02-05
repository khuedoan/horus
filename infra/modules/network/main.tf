locals {
  protocols = {
    all    = "all",
    icmp   = "1",
    icmpv6 = "58",
    tcp    = "6",
    udp    = "17"
  }
}

data "http" "public_ipv4" {
  url = "https://ipv4.icanhazip.com"
}

resource "oci_core_vcn" "vcn" {
  compartment_id = var.compartment_id
  display_name   = "vcn"
  cidr_blocks    = var.vcn_cidr_blocks
}

resource "oci_core_security_list" "base" {
  compartment_id = var.compartment_id
  display_name   = "base"
  vcn_id         = oci_core_vcn.vcn.id

  egress_security_rules {
    description = "Internet"
    protocol    = local.protocols["all"]
    destination = "0.0.0.0/0"
  }

  ingress_security_rules {
    description = "SSH"
    protocol    = local.protocols["tcp"]
    source      = "${chomp(data.http.public_ipv4.body)}/32"

    tcp_options {
      source_port_range {
        min = 1
        max = 65535
      }

      min = 22
      max = 22
    }
  }

  ingress_security_rules {
    description = "Kubernetes API"
    protocol    = local.protocols["tcp"]
    source      = "${chomp(data.http.public_ipv4.body)}/32"

    tcp_options {
      source_port_range {
        min = 1
        max = 65535
      }

      min = 6443
      max = 6443
    }
  }

  ingress_security_rules {
    description = "HTTP"
    protocol    = local.protocols["tcp"]
    source      = "0.0.0.0/0"

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

    tcp_options {
      source_port_range {
        min = 1
        max = 65535
      }

      min = 443
      max = 443
    }
  }

  ingress_security_rules {
    description = "SMTP"
    protocol    = local.protocols["tcp"]
    source      = "0.0.0.0/0"

    tcp_options {
      source_port_range {
        min = 1
        max = 65535
      }

      min = 25
      max = 25
    }
  }

  ingress_security_rules {
    description = "ESMTP Explicit TLS"
    protocol    = local.protocols["tcp"]
    source      = "0.0.0.0/0"

    tcp_options {
      source_port_range {
        min = 1
        max = 65535
      }

      min = 587
      max = 587
    }
  }

  ingress_security_rules {
    description = "ESMTP Implicit TLS"
    protocol    = local.protocols["tcp"]
    source      = "0.0.0.0/0"

    tcp_options {
      source_port_range {
        min = 1
        max = 65535
      }

      min = 465
      max = 465
    }
  }

  ingress_security_rules {
    description = "IMAP4 Implicit TLS"
    protocol    = local.protocols["tcp"]
    source      = "0.0.0.0/0"

    tcp_options {
      source_port_range {
        min = 1
        max = 65535
      }

      min = 993
      max = 993
    }
  }
}

resource "oci_core_subnet" "subnet" {
  compartment_id = var.compartment_id
  display_name   = "subnet"
  cidr_block     = var.subnet_cidr_block
  route_table_id = oci_core_vcn.vcn.default_route_table_id
  vcn_id         = oci_core_vcn.vcn.id

  security_list_ids = [
    oci_core_security_list.base.id
  ]
}

resource "oci_core_internet_gateway" "internet_gateway" {
  compartment_id = var.compartment_id
  display_name   = "internet-gateway"
  vcn_id         = oci_core_vcn.vcn.id
}

resource "oci_core_default_route_table" "default_route_table" {
  manage_default_resource_id = oci_core_vcn.vcn.default_route_table_id
  route_rules {
    destination       = "0.0.0.0/0"
    destination_type  = "CIDR_BLOCK"
    network_entity_id = oci_core_internet_gateway.internet_gateway.id
  }
}
