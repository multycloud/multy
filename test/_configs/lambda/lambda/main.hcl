multy "lambda" "super_long_function" {
  function_name = "super_long_function"
  runtime = "python3.9"
  source_code_dir = cloud_specific_value({aws: "source_dir/aws", azure: "source_dir/azure"})
}

config {
  location = var.location
  clouds   = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "EU_WEST_1"
}