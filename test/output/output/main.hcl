variable "var" {
  type    = string
  default = "$${join(\",\",[\"a\", \"b\"])}"
}

output "output1" {
  value = "string"
}

output "output3" {
  value = [1, "2", {"4": [5, 6]}]
}

output "output2" {
  value = 1
}

output "output4" {
  value = ["${var.var},c"]
}