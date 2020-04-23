package main

/*
  COMMENT
*/
//asdasd

import (
  "fmt"
  "os"
  "strconv"
)

func errorCheck(err error)  {
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }
}

func main()  {
  if os.Args[1] == "" || os.Args[2] == "" {
    fmt.Printf("Program start exanple 'go run NOD.go 556 356'\n")
  }

  bigger, err := strconv.Atoi(os.Args[1])
  errorCheck(err)

  smaller, err := strconv.Atoi(os.Args[2])
  errorCheck(err)

  result := 0

  for result == 0 {
    if bigger < smaller {
      bigger, smaller = smaller, bigger
    }

    if bigger % smaller != 0 {
      bigger = bigger % smaller
    } else {
      result = smaller
    }
  }

  fmt.Println("Result", result)
}
