variable "compartment_id" {
  description = "Compartment ID where to create all resources"
  type        = string
}

variable "shape" {
  description = "The shape configuration requested for the nodes"
  type = object({
    name   = string
    config = map(any)
  })
  validation {
    condition     = !(can(regex("Flex", var.shape.name)) && length(var.shape.config) == 0)
    error_message = "Shape config not found. Shape Config is required while using flexible shapes."
  }
}

variable "subnet_id" {
}

variable "ssh_public_key" {
  description = "SSH public key to add to all nodes"
}

variable "tags" {
  type    = map(string)
  default = {}
}

variable "image" {
  description = "OS image properties"

  type = object({
    operating_system = string
    version          = string
  })

  default = {
    operating_system = "Canonical Ubuntu"
    version          = "20.04"
  }
}
