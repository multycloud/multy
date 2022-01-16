config {
  location                    = var.location
  clouds = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "ireland"
}

multy virtual_network "example_vn_1" {
    name       = "example_vn"
    cidr_block = "1"
}

multy virtual_network "example_vn" {
    name       = example_vn_1.new_id
    cidr_block = max(1, 2)
}