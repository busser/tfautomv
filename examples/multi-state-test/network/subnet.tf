resource "terraform_data" "subnet" {
  input = {
    name    = "demo-network-subnet"
    purpose = "network"
  }
}
