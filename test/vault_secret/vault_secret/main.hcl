config {
  clouds = ["aws", "azure"]
  location = "uk"
}
multy "vault" "example" {
  name = "dev"
}
multy "vault_secret" "api_key" {
  name = "api_key"
  vault = example
  value = "xxx"
}