multy "lambda" "test" {
  function_name = "multyfunobj"
  runtime = "python3.9"
  source_code_object = source_code
}

multy "object_storage" "obj_storage" {
  name          = "test-storage-9999919"
  random_suffix = false
}

multy "object_storage_object" "source_code" {
  name                = "source_code.zip"
  object_storage      = obj_storage
  source              = cloud_specific_value({aws: "source_dir/aws_code.zip", azure: "source_dir/azure_code.zip"})
  acl                 = "public_read"
}

config {
  location = var.location
  clouds   = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "ireland"
}