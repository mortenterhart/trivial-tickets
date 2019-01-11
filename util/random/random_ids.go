package random

import (
    "math/rand"
    "time"
)

// letters are the valid characters for the ids
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// CreateRandomId generates a pseudo random id for tickets and mails
// Tweaked example from https://stackoverflow.com/a/22892986
func CreateRandomId(n int) string {

    // Seed the random function to make it more random
    rand.Seed(time.Now().UnixNano())

    // Create a slice, big enough to hold the id
    b := make([]rune, n)

    // Randomly choose a letter from the letters rune
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }

    return string(b)
}
