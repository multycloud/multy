multy "object_storage" "obj_storage1" {
  name = "mty-storage-001"
}
multy "object_storage" "obj_storage2" {
  name = "mty-storage-002"
}
multy "object_storage" "obj_storage3" {
  name = "mty-storage-003"
  // we're just testing that this doesn't break anything, in reality depends_on is not necessary here
  depends_on = [obj_storage1, obj_storage2]
}
config {
  location = "ireland"
  clouds   = ["aws", "azure"]
}