resource "terraform_data" "encryption" {
  input = {
    name    = "demo-shared-encryption"
    purpose = "shared"
  }
}
