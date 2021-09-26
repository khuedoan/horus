data "oci_identity_availability_domains" "availability_domains" {
  compartment_id = var.compartment_id
}

data "oci_core_images" "image" {
  compartment_id           = var.compartment_id
  operating_system         = var.image.operating_system
  operating_system_version = var.image.version
  shape                    = var.shape.name
  sort_by                  = "TIMECREATED"
  sort_order               = "DESC"
}

module "cloud_init" {
  source = "../cloud-init"
  role   = var.role
  token  = var.token
}

resource "oci_core_instance" "node" {
  compartment_id = var.compartment_id
  display_name   = "k3s-${var.role}"
  availability_domain = data.oci_identity_availability_domains.availability_domains.availability_domains[0].name
  freeform_tags  = var.tags

  create_vnic_details {
    assign_public_ip = true # TODO check if we can disable this
    subnet_id = var.subnet_id
  }

  extended_metadata = {
    subnet_id = var.subnet_id
  }

  metadata = {
    ssh_authorized_keys = var.ssh_public_key
    user_data = base64encode(module.cloud_init.cloud_init)
  }

  shape = var.shape.name
  source_details {
    source_type = "image"
    source_id   = data.oci_core_images.image.images[0].id
  }

  dynamic "shape_config" {
    for_each = length(var.shape.config) > 0 ? [1] : []
    content {
      ocpus         = tonumber(lookup(var.shape.config, "cpus", 0))
      memory_in_gbs = tonumber(lookup(var.shape.config, "memory", 0))
    }
  }
}
