resource "random_password" "token" {
  length  = 64
  special = false
}

resource "tls_private_key" "ssh" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P256"
}

resource "local_file" "ssh_private_key" {
  content         = tls_private_key.ssh.private_key_pem
  filename        = "${path.root}/private.pem"
  file_permission = "0600"
}


module "server_nodes" {
  count          = var.server_count
  source         = "../node"
  compartment_id = var.compartment_id
  subnet_id      = var.subnet_id
  ssh_public_key = tls_private_key.ssh.public_key_openssh
  role           = "server"
  token          = random_password.token.result
  shape = {
    name   = var.server_shape
    config = {
      cpus   = 2
      memory = 12
    }
  }
  tags = var.tags
}

module "agent_pool" {
  source         = "../node-pool"
  compartment_id = var.compartment_id
  subnet_id      = var.subnet_id
  ssh_public_key = tls_private_key.ssh.public_key_openssh
  role           = "agent"
  server_address = module.server_nodes[0].private_ip
  token          = random_password.token.result
  size           = var.agent_count
  shape = {
    name = var.agent_shape
    config = {
      cpus   = 2
      memory = 12
    }
  }
  tags = var.tags
}

resource "null_resource" "kubeconfig" {
  triggers = {
    # TODO optimize trigger
    file_exists = fileexists("${path.root}/kubeconfig.yaml")
  }

  provisioner "local-exec" {
    command = "ssh -o 'StrictHostKeyChecking no' -i ${local_file.ssh_private_key.filename} ubuntu@${module.server_nodes[0].public_ip} sudo cat /etc/rancher/k3s/k3s.yaml | sed 's/127.0.0.1/${module.server_nodes[0].public_ip}/g'> ${path.root}/kubeconfig.yaml"
  }
}
