resource "terraform_data" "service" {
  input = {
    name    = "demo-application-service"
    purpose = "application"
  }
}

resource "terraform_data" "worker" {
  input = {
    name    = "demo-application-worker"
    purpose = "application"
  }
}
