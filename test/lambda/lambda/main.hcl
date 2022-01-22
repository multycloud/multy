multy "lambda" "test" {
  function_name = "test_name"
  runtime = "python3.9"
  source_code_dir = cloud_specific_value({aws: "source_dir/aws", azure: "source_dir/azure"})
}

config {
  location = var.location
  clouds   = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "ireland"
}