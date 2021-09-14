variable "global_tags" {
  description = "Tags to add to all resources"
  type        = map(string)
  default = {
    project = "freecloud"
  }
}
