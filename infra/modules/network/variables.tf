variable "compartment_id" {
  default = "ocid1.compartment.oc1..aaaaaaaasvkl7yw6gj2pytybimo6ax7fg2jq7m4t5aueig3xrnfxwo7xwulq"
}

variable "vcn_cidr_blocks" {
  default = [
    "10.0.0.0/16"
  ]
}

variable "subnet_cidr_block" {
  default = "10.0.0.0/24"
}
