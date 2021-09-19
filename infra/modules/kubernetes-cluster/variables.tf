variable "compartment_id" {
  description = "Compartment ID where to create all resources"
  type        = string
}

variable "master_count" {
  description = "Number of master nodes"
}

variable "master_shape" {
  description = "The shape of master nodes"
  default     = "VM.Standard.E2.1.Micro"
}

variable "worker_count" {
  description = "Number of worker nodes"
}

variable "worker_shape" {
  description = "The shape of worker nodes"
  default = "VM.Standard.A1.Flex"
}

variable "subnet_id" {
}

variable "ssh_public_key" {
  description = "SSH public key to add to all nodes"
}

variable "tags" {
  type = map(string)
  default = {}
}
