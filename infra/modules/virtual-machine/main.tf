resource "oci_core_instance" "instance" {
  availability_domain = "gHLA:US-SANJOSE-1-AD-1" # TODO
  compartment_id      = var.compartment_id
  shape               = var.instance_shape

  instance_options {
    are_legacy_imds_endpoints_disabled = true
  }

  create_vnic_details {
    assign_public_ip = "true"
    subnet_id        = var.subnet_id
  }

  source_details {
    source_type = "image"
    source_id   = var.instance_image_id
  }

  metadata = {
    ssh_authorized_keys = var.ssh_public_key
    user_data           = filebase64("${path.module}/cloud-init.yaml")
  }
}
