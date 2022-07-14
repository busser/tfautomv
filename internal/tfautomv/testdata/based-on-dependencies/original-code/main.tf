resource "random_pet" "original_depender_first" {
    keepers = {
        other_pet = random_pet.original_dependee_first.id
    }
}

resource "random_pet" "original_dependee_first" {
    length = 1
}

resource "random_pet" "original_depender_second" {
    keepers = {
        other_pet = random_pet.original_dependee_second.id
    }
}

resource "random_pet" "original_dependee_second" {
    length = 2
}
