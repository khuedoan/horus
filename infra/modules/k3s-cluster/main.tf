module "server_pool" {
  source     = "../k3s-node-pool"
  node_count = var.server_count
}

module "agent_pool" {
  source     = "../k3s-node-pool"
  node_count = var.agent_count
}
