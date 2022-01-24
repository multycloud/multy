multy "lambda" "test" {
  function_name = "test_name"
  runtime = "python3.9"
  source_code_object = source_code
}

multy "object_storage" "obj_storage" {
  name          = "test-storage-9999919"
  random_suffix = false
}
multy "object_storage_object" "source_code" {
  name                = "source_code.zip"
  content             = "test"
  object_storage      = obj_storage
  content_type        = "application/zip"
}

config {
  location = var.location
  clouds   = ["azure", "aws"]
}
variable "location" {
  type    = string
  default = "ireland"
}