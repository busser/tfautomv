resource "terraform_data" "registry" {
  input = {
    name    = "demo-shared-registry"
    purpose = "move-between-states"
  }
}
