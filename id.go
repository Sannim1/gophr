package main

import (
    "fmt"
    "crypto/rand"
)

const idSource = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

const idSourceLength = byte(len(idSource))

func GenerateID(prefix string, length int) string {
    // create an array with the correct capacity
    id := make([]byte, length)

    // Fill our array with random numbers
    rand.Read(id)

    for i, b := range id {
        id[i] = idSource[b % idSourceLength]
    }

    return fmt.Sprintf("%s_%s", prefix, string(id))
}

// func main() {
//     var generatedId string
//     var idExists bool
//
//     duplicateFound := false
//     generatedIds := map[string]string{}
//
//     for ! duplicateFound {
//         generatedId = GenerateID("usr", 8)
//         fmt.Println(generatedId)
//
//         _, idExists = generatedIds[generatedId]
//
//         if idExists {
//             fmt.Println("found duplicate: %s", generatedId)
//             duplicateFound = true
//
//             return
//         }
//
//         generatedIds[generatedId] = generatedId
//     }
// }
