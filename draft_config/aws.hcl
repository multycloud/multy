variable "auth" {
  type = object({
    access_key        = string
    secret_access_key = string
  })
  default = {}
}

variable "cidr_block" {
  type = string
  default = "10.0.0.0/16"
}

