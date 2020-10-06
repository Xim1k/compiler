package main

import (
  "fmt"
  "regexp"
  "os"
)

func main()  {
  strings := []string{
    "string@asdasd.ru",
    "string@asd.ru",
    "string@a.ru",
  }

  subString := "asd"

  for _, string := range strings {
        matched, err := regexp.Match(subString, []byte(string))
        if matched {
            fmt.Printf("√ '%s' has subString\n", string)
        } else {
            fmt.Printf("X '%s' hasn't subString\n", string)
        }
    }
}
