resource "random_password" "token" {
  length  = 64
  special = false
}

module "server_nodes" {
  count          = var.server_count
  source         = "../single-node"
  compartment_id = var.compartment_id
  subnet_id      = var.subnet_id
  ssh_public_key = var.ssh_public_key
  role           = "server"
  token          = random_password.token.result
  shape = {
    name   = var.server_shape
    config = {}
  }
  tags = var.tags
}

module "agent_pool" {
  source         = "../node-pool"
  compartment_id = var.compartment_id
  subnet_id      = var.subnet_id
  ssh_public_key = var.ssh_public_key
  role           = "agent"
  token          = random_password.token.result
  size           = var.agent_count
  shape = {
    name = var.agent_shape
    config = {
      cpus = 2
      memory = 12
    }
  }
  tags = var.tags
}
