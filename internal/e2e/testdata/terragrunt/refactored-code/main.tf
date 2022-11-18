variable "prefix" {
  type = string
}

resource "random_pet" "refactored_first" {
  prefix = var.prefix
  length = 1
}

resource "random_pet" "refactored_second" {
  prefix = var.prefix
  length = 2
}

resource "random_pet" "refactored_third" {
  prefix = var.prefix
  length = 3
}
