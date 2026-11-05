resource "terraform_data" "configuration" {
  input = {
    name    = "demo-application-config"
    purpose = "application"
  }
}

output "configuration_id" {
  value = terraform_data.configuration.id
}
