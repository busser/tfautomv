resource "terraform_data" "vpc" {
  input = {
    name    = "demo-network-vpc"
    purpose = "network"
  }
}

output "vpc_id" {
  value = terraform_data.vpc.id
}
