# TODO: refactor time dependencies so we can test this
#multy "lambda" "function1" {
#  function_name = "privatemultyfun"
#  runtime = "python3.9"
#  source_code_object = private_source_code
#}

multy "lambda" "function2" {
  function_name      = "publicmultyfun"
  runtime            = "python3.9"
  source_code_object = public_source_code
}

multy "object_storage" "obj_storage" {
  name = "function-storage-1722"
}

#multy "object_storage_object" "private_source_code" {
#  name                = "source_code.zip"
#  object_storage      = obj_storage
#  source              = cloud_specific_value({aws: "source_dir/aws_code.zip", azure: "source_dir/azure_code.zip"})
#  acl                 = "private"
#}

multy "object_storage_object" "public_source_code" {
  name           = "source_code.zip"
  object_storage = obj_storage
  source         = cloud_specific_value({ aws : "source_dir/aws_code.zip", azure : "source_dir/azure_code.zip" })
  acl            = "public_read"
}


config {
  location = var.location
  clouds   = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "ireland"
}