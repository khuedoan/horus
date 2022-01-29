output "private_ip" {
  value = oci_core_instance.node.private_ip
}

output "public_ip" {
  value = oci_core_instance.node.public_ip
}
