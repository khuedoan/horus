variable "oracle_cloud" {
  description = "Oracle Cloud configuration"
  type = object({
    tenancy_ocid = string
    user_ocid    = string
    fingerprint  = string
    private_key  = string
    region       = string
  })
}

variable "compartment_name" {
  description = "Name of the compartment where to create all resources"
  type        = string
  default     = "horus"
}

variable "compartment_description" {
  description = "Description of the compartment where to create all resources"
  type        = string
  default     = "Horus Project"
}
