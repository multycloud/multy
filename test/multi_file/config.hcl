config {
  default_resource_group_name = "${resource_type}-rg"
  location                    = var.location
  clouds                      = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "ireland"
}