resource "terraform_data" "service" {
  input = {
    name    = "demo-application-service"
    purpose = "application"
  }
}
