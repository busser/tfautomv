resource "random_integer" "alpha" {
  min = 1
  max = 5
}

resource "random_integer" "beta" {
  min = 1
  max = 5
}

resource "random_pet" "first" {
  length    = random_integer.alpha.result
  separator = "-"
}

resource "random_pet" "second" {
  length    = random_integer.beta.result
  separator = "+"
}
