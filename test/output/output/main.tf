output "output1" {
  value = "string"
}
output "output2" {
  value = 1
}
output "output3" {
  value = [1, "2", { "4" = [5, 6] }]
}
output "output4" {
  value = ["${join(",", ["a", "b"])},c"]
}
