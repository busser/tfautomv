terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "3.3.2"
    }
  }
}

resource "random_pet" "original" {}

resource "random_id" "original" {
  byte_length = 6
}
