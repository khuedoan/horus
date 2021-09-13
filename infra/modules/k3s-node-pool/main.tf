module "virtual_machine" {
  count  = var.node_count
  source = "../virtual-machine"
}
