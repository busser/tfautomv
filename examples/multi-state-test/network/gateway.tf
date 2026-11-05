resource "terraform_data" "gateway" {
  input = {
    name    = "demo-network-gateway"
    purpose = "move-between-states"
  }
}

output "gateway_id" {
  value = terraform_data.gateway.id
}
