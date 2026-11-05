resource "terraform_data" "vpc" {
  input = {
    name    = "demo-network-vpc"
    purpose = "network"
  }
}
