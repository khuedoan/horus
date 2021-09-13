variable "instance_shape" {
  default = "VM.Standard.E2.1.Micro"
}

variable "compartment_id" {
  default = "ocid1.compartment.oc1..aaaaaaaasvkl7yw6gj2pytybimo6ax7fg2jq7m4t5aueig3xrnfxwo7xwulq"
}

variable "instance_image_id" {
  default = "ocid1.image.oc1.us-sanjose-1.aaaaaaaan4g4q527bljtyczck6xrsutbzps6h7mut2xcfhnbzw66sbbsvwoq"
}

variable "vcn_cidr_blocks" {
  default = [
    "10.0.0.0/16"
  ]
}

variable "subnet_cidr_block" {
  default = "10.0.0.0/24"
}
