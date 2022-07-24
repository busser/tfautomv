resource "random_integer" "first" {
  min = 1
  max = 5
}

resource "random_integer" "second" {
  min = 1
  max = 5
}

resource "random_pet" "first" {
  length    = random_integer.first.result
  separator = "-"
}

resource "random_pet" "second" {
  length    = random_integer.second.result
  separator = "+"
}
