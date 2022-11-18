variable "prefix" {
  type = string
}

resource "random_pet" "original_first" {
  prefix = var.prefix
  length = 1
}

resource "random_pet" "original_second" {
  prefix = var.prefix
  length = 2
}

resource "random_pet" "original_third" {
  prefix = var.prefix
  length = 3
}
