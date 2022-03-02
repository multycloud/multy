multy "virtual_network" "example_vn" {
  name       = "example_vn"
  cidr_block = "10.0.0.0/16"
}
config {
  location = var.location
  clouds   = ["aws"]
}
variable "location" {
  type    = string
  default = "ireland"
}

output "test_nested" {
  value = file(file("file.txt"))
}
output "test_var" {
  value = file("${var.location}.txt")
}
output "test_output_var" {
  value = "${file("${aws.example_vn.id}.txt")}\n"
}
output "test_quotes" {
  value = "${file("\"")}"
}