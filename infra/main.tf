module "k3s_cluster" {
  source       = "./modules/k3s-cluster"
  server_count = 1
  agent_count  = 3
}
