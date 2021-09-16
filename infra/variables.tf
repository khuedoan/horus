variable "compartment_id" {
  description = "Compartment ID where to create all resources"
  type        = string
}

variable "common_tags" {
  description = "Common tags for all resources"
  type        = map(string)
  default = {
    project = "freecloud"
  }
}
