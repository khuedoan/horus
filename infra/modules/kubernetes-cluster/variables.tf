variable "compartment_id" {
  description = "Compartment ID where to create all resources"
  type        = string
}

variable "server_count" {
  description = "Number of server nodes"
}

variable "server_shape" {
  description = "The shape of server nodes"
  default     = "VM.Standard.A1.Flex"
}

variable "agent_count" {
  description = "Number of agent nodes"
}

variable "agent_shape" {
  description = "The shape of agent nodes"
  default     = "VM.Standard.A1.Flex"
}

variable "subnet_id" {
}

variable "tags" {
  type    = map(string)
  default = {}
}
