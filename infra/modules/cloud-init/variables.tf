variable "cluster_init" {
  description = "Initial server node in the cluster"
  type        = bool
  default     = false
}

variable "role" {
  description = "Node role"
  type        = string
  validation {
    condition     = contains(["server", "agent"], var.role)
    error_message = "Node role must be server or agent."
  }
}

variable "token" {
  description = "Shared secret used to join a server or agent to a cluster"
  type        = string
  sensitive   = true
}

variable "server_address" {
  type = string
}
