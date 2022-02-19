multy "lambda" "lambda" {
  function_name = "multy_function"
  runtime = "python3.9"
  source_code_dir = cloud_specific_value({aws: "source_dir/aws", azure: "source_dir/azure"})
}

multy "object_storage" "obj_storage" {
  clouds = ["aws"]
  name          = "test-storage-9999919"
  random_suffix = false
}
multy "object_storage_object" "file1_public" {
  clouds = ["aws"]
  name                = "index.html"
  content             = "<button onclick='lambda(\"${aws.lambda.url}\")'>Call aws</button><button onclick='lambda(\"${azure.lambda.url}/trigger\")'>Call azure</button><script>function lambda(url) {fetch(url).then(data => data.text()).then(data => alert(data));}</script>"
  object_storage      = obj_storage
  content_type        = "text/html"
  acl = "public_read"
}

config {
  location = var.location
  clouds   = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "ireland"
}