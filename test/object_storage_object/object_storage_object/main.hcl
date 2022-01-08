config {
  default_resource_group_name = "${resource_type}-rg"
  location                    = "ireland"
  clouds                      = ["aws", "azure"]
}
multy "object_storage" "obj_storage" {
  name          = "test-storage-9999919"
  random_suffix = false
}
multy "object_storage_object" "file1" {
  name                = "index.html"
  content             = "<h1>Hi</h1>"
  object_storage_name = obj_storage.id
  content_type        = "text/html"
}