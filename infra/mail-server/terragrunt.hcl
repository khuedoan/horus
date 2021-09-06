include {
  path = find_in_parent_folders("config.hcl")
}

terraform {
  source = "../modules//vm"
}

inputs = {
}
