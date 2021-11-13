resource "oci_identity_dynamic_group" "kubernetes_nodes" {
  compartment_id = var.tenancy_id
  name           = "kubernetes-nodes"
  description    = "Dynamic group for Kubernetes Cloud Controller Manager"
  # TODO reduce matching rule
  matching_rule = "All {instance.compartment.id = '${var.compartment_id}'}"
}

resource "oci_identity_policy" "cloud_controller_manager" {
  compartment_id = var.compartment_id
  name           = "cloud-controller-manager"
  description    = "Policy to allow Kubernetes nodes to call services"
  statements     = [
    # https://github.com/oracle/oci-cloud-controller-manager/blob/master/manifests/provider-config-example.yaml
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to read instance-family in compartment id ${var.compartment_id}",
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to use virtual-network-family in compartment id ${var.compartment_id}",
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to manage load-balancers in compartment id ${var.compartment_id}",

    # https://github.com/oracle/oci-cloud-controller-manager/blob/master/flex-volume-provisioner.md
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to manage volumes in compartment id ${var.compartment_id}",
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to manage file-systems in compartment id ${var.compartment_id}",

    # https://github.com/oracle/oci-cloud-controller-manager/blob/master/flex-volume-driver.md
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to read vnic-attachments in compartment id ${var.compartment_id}",
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to read vnics in compartment id ${var.compartment_id}",
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to read instances in compartment id ${var.compartment_id}",
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to read subnets in compartment id ${var.compartment_id}",
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to use volumes in compartment id ${var.compartment_id}",
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to use instances in compartment id ${var.compartment_id}",
    "Allow dynamic-group ${oci_identity_dynamic_group.kubernetes_nodes.name} to manage volume-attachments in compartment id ${var.compartment_id}",
  ]
}
