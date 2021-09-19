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

variable "compartment_id" {
  default = "ocid1.compartment.oc1..aaaaaaaasvkl7yw6gj2pytybimo6ax7fg2jq7m4t5aueig3xrnfxwo7xwulq"
}

