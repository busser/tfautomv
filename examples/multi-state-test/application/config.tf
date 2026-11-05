resource "terraform_data" "configuration" {
  input = {
    name    = "demo-application-config"
    purpose = "application"
  }
}
