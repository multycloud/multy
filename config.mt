config {
  location                    = var.location
  clouds = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "ireland"
}
multy virtual_network "example_vn" {
  name       = "example_vn"
  cidr_block = ""
}