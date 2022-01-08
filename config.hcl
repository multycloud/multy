multy "lambda" "test" {
  function_name = "test_name"
  os = "linux"
  runtime = "nodejs12.x"
  source_code_dir = "dir"
}

config {
  location                    = var.location
  clouds = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "ireland"
}