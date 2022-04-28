multy "virtual_network" "example_vn" {
  name       = "example_vn"
  cidr_block = cloud_specific_value({ aws : "10.0.0.0/24", default : "10.0.0.1/24" })
}
config {
  location = var.location
  clouds   = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "EU_WEST_1"
}