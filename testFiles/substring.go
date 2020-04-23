package main

import (
  "fmt"
  "regexp"
  "os"
)

func errorCheck(err error)  {
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }
}

func main()  {
  strings := []string{
    "string@asdasd.ru",
    "string@asd.ru",
    "string@a.ru",
  }

  subString := "asd"

  for _, string := range strings {
        matched, err := regexp.Match(subString, []byte(string))
        errorCheck(err)
        if matched {
            fmt.Printf("√ '%s' has subString\n", string)
        } else {
            fmt.Printf("X '%s' hasn't subString\n", string)
        }
    }
}
