module "worker" {
  source = "../node-pool"
  compartment_id = var.compartment_id
  subnet_id = var.subnet_id
  ssh_public_key = var.ssh_public_key
  size = var.worker_count
  shape = var.worker_shape
  shape_config = {
    cpus = 2
    memory = 12
  }
  tags = var.tags
}
