package main

import "fmt"

func encodeMsg(rawMsg string) string{
    var out string
    user := GetUser()
    out = fmt.Sprintf("%s:  %s", user, rawMsg)
    return out
}
