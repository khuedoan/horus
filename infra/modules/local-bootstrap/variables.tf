variable "cluster" {
  type = string
}

variable "credentials" {
  type = object({
    client_certificate     = string
    client_key             = string
    cluster_ca_certificate = string
    host                   = string
  })
}

variable "cluster_domain" {
  type    = string
  default = "127-0-0-1.nip.io"
}
