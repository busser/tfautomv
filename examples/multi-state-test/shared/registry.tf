resource "terraform_data" "registry" {
  input = {
    name    = "demo-shared-registry"
    purpose = "move-between-states"
  }
}

output "registry_id" {
  value = terraform_data.registry.id
}
