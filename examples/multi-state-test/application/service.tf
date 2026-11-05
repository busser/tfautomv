resource "terraform_data" "service" {
  input = {
    name    = "demo-application-service"
    purpose = "application"
  }
}

output "service_id" {
  value = terraform_data.service.id
}
