variable "compartment_id" {
  description = "Compartment ID where to create all resources"
  type        = string
}

variable "shape" {
  description = "The shape of the nodes"
}

variable "shape_config" {
  description = "The shape configuration requested for the nodes"
  type = object({
    cpus   = number
    memory = number
  })
  default = {
    cpus   = null
    memory = null
  }

  # TODO https://github.com/hashicorp/terraform/issues/25609
  # validation {
  #   condition = !(can(regex("Flex", var.image_id)) && length(var.shape_config) == 0)
  #   error_message = "Shape config not found. Shape Config is required while using flexible shapes."
  # }
}

variable "size" {
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
