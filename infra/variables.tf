variable "parent_compartment_id" {
  description = "Compartment ID where to create the subcompartment for this project (usually the root compartment)"
  type        = string
}

variable "compartment_name" {
  description = "Name of the compartment where to create all resources"
  type        = string
  default     = "freecloud"
}

variable "compartment_description" {
  description = "Description of the compartment where to create all resources"
  type        = string
  default     = "freecloud"
}

variable "common_tags" {
  description = "Common tags for all resources"
  type        = map(string)
  default = {
    project = "freecloud"
  }
}
