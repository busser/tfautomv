resource "terraform_data" "worker" {
  input = {
    name    = "demo-application-worker"
    purpose = "application"
  }
}
