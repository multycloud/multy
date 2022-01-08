config {
  location                    = "uk"
  clouds                      = ["azure"]
  default_resource_group_name = "${resource_type}-${rg_vars.app}-rg"
}
multy "virtual_network" "example_vn" {
  rg_vars    = {
    app = "backend"
  }
  name       = "example_vn"
  cidr_block = "10.0.0.0/16"
}