variable "compartment_id" {
  description = "Compartment ID where to create all resources"
  type        = string
}

variable "server_count" {
  description = "Number of server (master) nodes"
}

variable "agent_count" {
  description = "Number of agent (worker) nodes"
}

variable "subnet_id" {
}

variable "ssh_public_key" {
  description = "SSH public key to add to all nodes"
}

variable "image_operating_system" {
  default = "Canonical Ubuntu"
}

variable "image_operating_system_version" {
  default = "20.04"
}
