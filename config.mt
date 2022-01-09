config {
  location                    = var.location
  clouds = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "ireland"
}