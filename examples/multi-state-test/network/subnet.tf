resource "terraform_data" "subnet" {
  input = {
    name    = "demo-network-subnet"
    purpose = "network"
  }
}

output "subnet_id" {
  value = terraform_data.subnet.id
}
