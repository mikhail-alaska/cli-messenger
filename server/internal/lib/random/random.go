package random

import "math/rand"

func NewRandomString(length int) string{
    out := ""
    alphabet := "1234567890-_+qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
    for i := 0; i< length; i++{
        temp := rand.Intn(len(alphabet))
        out+= string(alphabet[temp])
    }
    return out
}
