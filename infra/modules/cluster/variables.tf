variable "vault_password" {
  type      = string
  sensitive = true
}

variable "instance_public_ip" {
  type      = string
  sensitive = true
}

variable "ssh_private_key" {
  type = string
}
