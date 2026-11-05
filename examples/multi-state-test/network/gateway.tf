resource "terraform_data" "gateway" {
  input = {
    name    = "demo-network-gateway"
    purpose = "move-between-states"
  }
}
