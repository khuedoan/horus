variable "tenancy_id" {
  description = "Tenancy ID where to create all resources"
  type        = string
}

variable "compartment_id" {
  description = "Compartment ID where to create all resources"
  type        = string
}

variable "vcn_id" {
  description = "VCN ID where to create all resources"
  type        = string
}

variable "subnet_id" {
  description = "Subnet ID where to create all resources"
  type        = string
}
