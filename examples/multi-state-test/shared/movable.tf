resource "terraform_data" "movable" {
  input = {
    name    = "demo-shared-movable-resource"
    purpose = "move-between-states"
  }
}
