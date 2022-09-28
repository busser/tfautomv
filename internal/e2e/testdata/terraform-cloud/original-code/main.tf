terraform {
  cloud {
    organization = "busser"

    workspaces {
      name = "tfautomv"
    }
  }
}

resource "random_pet" "original_first" {
  length = 1
}

resource "random_pet" "original_second" {
  length = 2
}

resource "random_pet" "original_third" {
  length = 3
}
