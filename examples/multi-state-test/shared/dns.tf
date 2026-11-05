resource "terraform_data" "dns" {
  input = {
    name    = "demo-shared-dns"
    purpose = "shared"
  }
}
