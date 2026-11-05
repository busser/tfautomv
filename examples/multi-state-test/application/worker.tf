resource "terraform_data" "worker" {
  input = {
    name    = "demo-application-worker"
    purpose = "application"
  }
}

output "worker_id" {
  value = terraform_data.worker.id
}
