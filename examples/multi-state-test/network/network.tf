resource "terraform_data" "vpc" {
  input = {
    name    = "demo-network-vpc"
    purpose = "network"
  }
}

resource "terraform_data" "subnet" {
  input = {
    name    = "demo-network-subnet"
    purpose = "network"
  }
}
