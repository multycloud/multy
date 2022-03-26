config {
  default_resource_group_name = "${resource_type}-rg"
  location                    = "ireland"
  clouds                      = ["aws", "azure"]
}
multy "object_storage" "obj_storage" {
  name = "test-storage-9999919"
}
multy "object_storage_object" "file1_public" {
  name           = "index.html"
  content        = cloud_specific_value({ aws : "<h1>Hi from AWS</h1>", azure : "<h1>Hi from Azure</h1>" })
  object_storage = obj_storage
  content_type   = "text/html"
  acl            = "public_read"
}
multy "object_storage_object" "file2_private" {
  name           = "index_private.html"
  content        = "<h1>Hi</h1>"
  object_storage = obj_storage
  content_type   = "text/html"
}
multy "object_storage_object" "file3_source" {
  name           = "index.html"
  source         = "test.zip"
  object_storage = obj_storage
  acl            = "public_read"
}