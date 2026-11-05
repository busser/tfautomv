resource "terraform_data" "movable" {
  input = {
    name    = "demo-movable-resource"
    purpose = "move-between-states"
  }
}
