config {
  clouds = ["aws", "azure"]
  location = "uk"
}
multy "vault" "example" {
  name = "dev-test-secret-multy"
}
multy "vault_secret" "api_key" {
  name = "api-key"
  vault = example
  value = "xxx"
}