package main

import (
  "fmt"
  "strings"
)

func main() {
    var stringss = [3]string{"string@asdasd.ru","string@asd.ru","string@a.ru"}
    var index int = 0

    var subString string = "asd"

    for index < len(stringss) {
        if strings.Contains(stringss[index], subString) {
            fmt.Printf("√ '%s' has subString\n", stringss[index])
        } else {
            fmt.Printf("X '%s' hasn't subString\n", stringss[index])
        }
    }
}
