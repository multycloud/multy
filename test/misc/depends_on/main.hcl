multy "object_storage" "obj_storage1" {
  name = "mty-storage-001"
}
multy "object_storage" "obj_storage2" {
  name = "mty-storage-002"
}
multy "object_storage" "obj_storage3" {
  name = "mty-storage-003"
  depends_on = [obj_storage1, obj_storage2]
}
config {
  location = "ireland"
  clouds   = ["aws", "azure"]
}