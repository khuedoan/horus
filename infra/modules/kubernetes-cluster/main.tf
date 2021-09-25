resource "random_password" "token" {
  length           = 64
  special          = false
}

module "master_pool" {
  source = "../node-pool"
  compartment_id = var.compartment_id
  subnet_id = var.subnet_id
  ssh_public_key = var.ssh_public_key
  role = "server"
  token = random_password.token.result
  size = var.master_count
  shape = {
    name = var.master_shape
    config = {}
  }
  tags = var.tags
}

# TODO workaround until there's ARM capacity
# module "worker_pool" {
#   source = "../node-pool"
#   compartment_id = var.compartment_id
#   subnet_id = var.subnet_id
#   ssh_public_key = var.ssh_public_key
#   size = var.worker_count
#   shape = {
#     name = var.worker_shape
#     config = {
#       cpus = 2
#       memory = 12
#     }
#   }
#   tags = var.tags
# }

module "worker_pool_temp" {
  source = "../node-pool"
  compartment_id = var.compartment_id
  subnet_id = var.subnet_id
  ssh_public_key = var.ssh_public_key
  role = "agent"
  token = random_password.token.result
  size = var.master_count
  shape = {
    name = var.master_shape
    config = {}
  }
  tags = var.tags
}
