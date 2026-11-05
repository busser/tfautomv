resource "terraform_data" "encryption" {
  input = {
    name    = "demo-shared-encryption"
    purpose = "shared"
  }
}

output "encryption_id" {
  value = terraform_data.encryption.id
}
