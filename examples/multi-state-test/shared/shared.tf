resource "terraform_data" "dns" {
  input = {
    name    = "demo-shared-dns"
    purpose = "shared"
  }
}

resource "terraform_data" "encryption" {
  input = {
    name    = "demo-shared-encryption"
    purpose = "shared"
  }
}
